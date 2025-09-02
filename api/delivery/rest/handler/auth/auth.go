package auth

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
}

type authHandler struct {
	authUsecase authDom.Usecase
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	input, err := decodeLoginInput(c)
	if err != nil {
		return err
	}

	output, err := h.authUsecase.Login(c.UserContext(), input)
	if err != nil {
		return err
	}
	c.Cookie(buildAuthCookie(output))
	return c.Status(fiber.StatusOK).JSON(response.Success("login successful", output))
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	input, err := decodeRegisterInput(c)
	if err != nil {
		return err
	}

	session, err := h.authUsecase.Register(c.UserContext(), input)
	if err != nil {
		return err
	}
	if session == nil {
		return c.Status(fiber.StatusAccepted).JSON(response.Success[any]("check your mail for verification", nil))
	}
	c.Cookie(buildAuthCookie(session))
	return c.Status(fiber.StatusOK).JSON(response.Success("registration successful", session))
}

func buildAuthCookie(session *authDom.Session) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     utils.AuthTokenCookieName,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		SameSite: fiber.CookieSameSiteStrictMode,
		HTTPOnly: true,
	}
}

func New(authUsecase authDom.Usecase) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}
