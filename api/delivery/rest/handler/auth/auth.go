package auth

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/handler"
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type authHandler struct {
	authUc authDom.Usecase
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var input authDom.LoginInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	output, err := h.authUc.Login(c.UserContext(), &input)
	if err != nil {
		return err
	}

	c.Cookie(response.BuildAuthCookie(output))

	return c.Status(fiber.StatusOK).JSON(response.Success("login successful", output))
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	var input authDom.RegisterInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	session, err := h.authUc.Register(c.UserContext(), &input)
	if err != nil {
		return err
	}
	if session == nil {
		return c.Status(fiber.StatusAccepted).JSON(response.Success[any]("check your mail for verification", nil))
	}

	c.Cookie(response.BuildAuthCookie(session))

	return c.Status(fiber.StatusOK).JSON(response.Success("registration successful", session))
}

func (h *authHandler) Verify(c *fiber.Ctx) error {
	token := c.Query(utils.TokenQueryKey)
	if token == "" {
		return errorDom.Raise(c.UserContext(), errorDom.ErrRestMissingToken, "", nil, nil)
	}

	if err := h.authUc.Verify(c.UserContext(), token); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any]("verified successfully! proceed to login", nil))
}

func (h *authHandler) ForgotPassword(c *fiber.Ctx) error {
	var input authDom.ForgotPasswordInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	if err := h.authUc.ForgotPassword(c.UserContext(), &input); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any]("check your mail for password reset", nil))
}

func (h *authHandler) ResetPassword(c *fiber.Ctx) error {
	var input authDom.ResetPasswordInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	if err := h.authUc.ResetPassword(c.UserContext(), &input); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any]("password reset successful", nil))
}

func (h *authHandler) Logout(c *fiber.Ctx) error {
	if err := h.authUc.Logout(c.UserContext()); err != nil {
		return err
	}

	c.Cookie(response.BuildAuthCookie(&authDom.Session{}))

	return c.SendStatus(fiber.StatusNoContent)
}

func New(authUsecase authDom.Usecase) AuthHandler {
	return &authHandler{
		authUc: authUsecase,
	}
}
