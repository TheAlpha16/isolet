package email

import (
	"context"
	"embed"
	"sync"

	"github.com/TheAlpha16/isolet/api/infra/smtp"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	emailDom "github.com/TheAlpha16/isolet/api/internal/domain/email"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"go.uber.org/zap"
)

//go:embed templates/*
var templateFS embed.FS

type emailImpl struct {
	client smtp.SMTP
}

func (e *emailImpl) SendEmailAsync(ctx context.Context, input *emailDom.EmailInput) error {
	return nil
}

func New(ctx context.Context, wg *sync.WaitGroup, cvUc cvDom.Usecase) emailDom.Usecase {
	config := utils.GetConfig()

	client, err := smtp.NewSMTP(ctx, &smtp.Config{
		Creds: smtp.Creds{
			Host:     cvUc.GetString(ctx, cvDom.SMTPHost),
			Port:     cvUc.GetInt(ctx, cvDom.SMTPPort),
			User:     cvUc.GetString(ctx, cvDom.SMTPUser),
			Password: cvUc.GetString(ctx, cvDom.SMTPPassword),
		},
		Retries:     config.SMTP.Retries,
		ConnTimeout: config.SMTP.Timeout,
	}, wg)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		logger.GetAppLogger().Fatal("failed to initialize email service", zap.Error(err))
	}

	return &emailImpl{client: client}
}
