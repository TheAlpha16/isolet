package email

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"net/url"
	"sync"

	"github.com/TheAlpha16/isolet/api/infra/smtp"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
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
	cvUc   cvDom.Usecase
}

type config struct {
	Subject      string
	Template     string
	RedirectPath string
}

var configs = map[emailDom.Type]config{
	emailDom.TypeVerification: {Subject: "Verify your email", Template: "templates/verification.html", RedirectPath: utils.RouteAuthVerify},
	// emailDom.TypePasswordReset: {Subject: "Password Reset", Template: "templates/password_reset.html", RedirectPath: "/reset"},
}

func (e *emailImpl) SendEmailAsync(ctx context.Context, input *emailDom.EmailInput) error {
	config, ok := configs[input.Type]
	if !ok {
		return errorDom.Raise(ctx, errorDom.ErrEmailInvalidType, "", nil, common.ExtraData{"type": input.Type})
	}

	link, err := e.buildLink(ctx, config.RedirectPath, input.Token)
	if err != nil {
		return err
	}

	body, err := e.getBody(ctx, config.Template, &emailDom.TemplateInput{
		EventName: e.cvUc.GetString(ctx, cvDom.EventName),
		Username:  input.Username,
		Link:      link,
	})
	if err != nil {
		return err
	}

	if err := e.client.SendAsync(ctx, &emailDom.Email{
		Sender: emailDom.Entity{
			Name:    e.cvUc.GetString(ctx, cvDom.EventName),
			Address: e.cvUc.GetString(ctx, cvDom.EmailSender),
		},
		Recipients: []emailDom.Entity{
			{
				Name:    input.Username,
				Address: input.To,
			},
		},
		Subject: config.Subject,
		Body:    body,
	}); err != nil {
		return err
	}

	return nil
}

func (e *emailImpl) getBody(ctx context.Context, templatePath string, data *emailDom.TemplateInput) (string, error) {
	template, err := template.ParseFS(templateFS, templatePath)
	if err != nil {
		return "", errorDom.Raise(ctx, errorDom.ErrEmailTemplateFetch, "", err, common.ExtraData{"path": templatePath})
	}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", errorDom.Raise(ctx, errorDom.ErrEmailTemplateExecute, "", err, common.ExtraData{"path": templatePath})
	}

	return buf.String(), nil
}

func (e *emailImpl) buildLink(ctx context.Context, redirectPath string, token string) (string, error) {
	config := utils.GetConfig()
	publicURL := e.cvUc.GetString(ctx, cvDom.PublicURL)

	url, err := url.Parse(publicURL)
	if err != nil {
		return "", errorDom.Raise(ctx, errorDom.ErrEmailLinkBuild, "failed to parse public URL", err, common.ExtraData{"public_url": publicURL})
	}

	url = url.JoinPath(config.Rest.APIVersionPrefix, redirectPath)
	query := url.Query()
	query.Set(utils.TokenQueryKey, token)
	url.RawQuery = query.Encode()

	if url.Scheme == "" {
		url.Scheme = "https"
	}

	return url.String(), nil
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

	return &emailImpl{
		client: client,
		cvUc:   cvUc,
	}
}
