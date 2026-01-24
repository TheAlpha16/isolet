package smtp

import (
	"context"
	"sync"
	"time"

	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	emailDom "github.com/TheAlpha16/isolet/oracle/internal/domain/email"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"
	"github.com/TheAlpha16/isolet/oracle/utils/tracer"

	"github.com/go-gomail/gomail"
	"go.uber.org/zap"
)

type smtpImpl struct {
	config   *Config
	sender   gomail.SendCloser
	messages chan *gomail.Message
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	mu       sync.Mutex
}

func (s *smtpImpl) SendAsync(ctx context.Context, email *emailDom.Email) error {
	message := gomail.NewMessage()
	recipients := make([]string, len(email.Recipients))

	for i, recipient := range email.Recipients {
		recipients[i] = message.FormatAddress(recipient.Address, recipient.Name)
	}

	message.SetAddressHeader("From", email.Sender.Address, email.Sender.Name)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", email.Subject)
	message.SetBody("text/html", email.Body)

	select {
	case <-s.ctx.Done():
		return errorDom.Raise(ctx, errorDom.ErrSMTPContextCanceled, "", nil, common.ExtraData{"recipients": recipients})
	case s.messages <- message:
		return nil
	default:
		return errorDom.Raise(ctx, errorDom.ErrSMTPQueueFull, "", nil, common.ExtraData{"recipients": recipients})
	}
}

func (s *smtpImpl) getSender() (gomail.SendCloser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	if s.sender == nil {
		dialer := gomail.NewDialer(
			s.config.Host, s.config.Port, s.config.User, s.config.Password,
		)
		s.sender, err = dialer.Dial()
		if err != nil {
			return nil, errorDom.Raise(s.ctx, errorDom.ErrSMTPDialError, "", err, nil)
		}
	}
	return s.sender, nil
}

func (s *smtpImpl) sendEmail(message *gomail.Message) error {
	sender, err := s.getSender()
	if err != nil {
		return err
	}
	err = gomail.Send(sender, message)
	if err != nil {
		return errorDom.Raise(s.ctx, errorDom.ErrSMTPSendFailed, "", err, common.ExtraData{"recipients": message.GetHeader("To"), "subject": message.GetHeader("Subject")})
	}
	return nil
}

func (s *smtpImpl) handleEmail(message *gomail.Message) {
	var err error

	for i := 0; i <= s.config.Retries; i++ {
		if err = s.sendEmail(message); err == nil {
			return
		}
		logger.GetAppLogger().Error("error sending email", zap.Error(err))
		tracer.Sleep(s.ctx, "smtp.handleEmail.retry", time.Second*time.Duration(1<<i)) // 1s, 2s, 4s, 8s...
	}

	logger.GetAppLogger().Error("email max retries exceeded", zap.String("recipient", message.GetHeader("To")[0]))
	errorDom.RaiseToSentry(s.ctx, err)
}

func (s *smtpImpl) start() error {
	// DEBUG
	// if _, err := s.getSender(); err != nil {
	// 	return err
	// }

	go func() {
		resetChan := make(chan struct{}, 1)
		timer := time.NewTimer(s.config.ConnTimeout)
		defer timer.Stop()

		for {
			select {
			case message := <-s.messages:
				s.wg.Add(1)
				go func() {
					defer s.wg.Done()
					s.handleEmail(message)

					select {
					case resetChan <- struct{}{}:
					default:
					}
				}()
			case <-resetChan:
				if !timer.Stop() {
					<-timer.C // drain channel if already expired
				}
				timer.Reset(s.config.ConnTimeout)
			case <-timer.C:
				s.closeSender()
				timer.Reset(s.config.ConnTimeout)
			case <-s.ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (s *smtpImpl) closeSender() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sender != nil {
		if err := s.sender.Close(); err != nil {
			errorDom.RaiseToSentry(s.ctx, err)
			logger.GetAppLogger().Error("failed to close SMTP sender", zap.Error(err))
		}
	}
	s.sender = nil
}

func NewSMTP(ctx context.Context, config *Config, wg *sync.WaitGroup) (SMTP, error) {
	ctx, cancel := context.WithCancel(ctx)

	smtp := &smtpImpl{
		config:   config,
		messages: make(chan *gomail.Message, config.ChannelSize),
		ctx:      ctx,
		cancel:   cancel,
		wg:       wg,
	}

	if err := smtp.start(); err != nil {
		return nil, err
	}

	utils.InterruptHandlerChannel <- func() {
		cancel()
	}

	return smtp, nil
}
