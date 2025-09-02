package auth

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"

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
	return c.Status(fiber.StatusOK).JSON(response.Success("login successful", output))
}

func (h *authHandler) Register(c *fiber.Ctx) error {
	input, err := decodeRegisterInput(c)
	if err != nil {
		return err
	}

	err = h.authUsecase.Register(c.UserContext(), input)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(response.Success[any]("registration successful", nil))
}

func New(authUsecase authDom.Usecase) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}
