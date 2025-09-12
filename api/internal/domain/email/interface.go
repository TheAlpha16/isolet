package email

import "context"

type Usecase interface {
	SendEmailAsync(ctx context.Context, input *EmailInput) error
}
