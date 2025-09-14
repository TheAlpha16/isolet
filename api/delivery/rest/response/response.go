package response

import (
	"time"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type APIResponse[T any] struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func Success[T any](message string, data T) APIResponse[T] {
	return newAPIResponse(StatusSuccess, message, data)
}

func Error[T any](message string, data T) APIResponse[T] {
	return newAPIResponse(StatusError, message, data)
}

func newAPIResponse[T any](status Status, message string, data T) APIResponse[T] {
	return APIResponse[T]{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func BuildAuthCookie(session *authDom.Session) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     utils.AuthTokenCookieName,
		Value:    session.Token,
		Expires:  time.Unix(session.ExpiresAt, 0),
		SameSite: fiber.CookieSameSiteStrictMode,
		HTTPOnly: true,
	}
}
