package team

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

type TeamHandler interface {
	Create(c *fiber.Ctx) error
	Join(c *fiber.Ctx) error
	GenerateInvite(c *fiber.Ctx) error
	AcceptInvite(c *fiber.Ctx) error
}

type teamHandler struct {
	teamUc teamDom.Usecase
}

func (h *teamHandler) Create(c *fiber.Ctx) error {
	input, err := decodeCreateRequest(c)
	if err != nil {
		return err
	}

	session, err := h.teamUc.Create(c.UserContext(), input)
	if err != nil {
		return err
	}
	c.Cookie(response.BuildAuthCookie(session))
	return c.Status(fiber.StatusCreated).JSON(response.Success("team created successfully", session))
}

func (h *teamHandler) Join(c *fiber.Ctx) error {
	input, err := decodeJoinRequest(c)
	if err != nil {
		return err
	}

	session, err := h.teamUc.Join(c.UserContext(), input)
	if err != nil {
		return err
	}
	c.Cookie(response.BuildAuthCookie(session))
	return c.Status(fiber.StatusCreated).JSON(response.Success("team joined successfully", session))
}

func (h *teamHandler) GenerateInvite(c *fiber.Ctx) error {
	output, err := h.teamUc.GenerateInvite(c.UserContext())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success("invite link generated successfully", output))
}

func (h *teamHandler) AcceptInvite(c *fiber.Ctx) error {
	token := c.Query(utils.TokenQueryKey)
	if token == "" {
		return errorDom.Raise(c.UserContext(), errorDom.ErrRestMissingToken, "", nil, nil)
	}

	session, err := h.teamUc.AcceptInvite(c.UserContext(), token)
	if err != nil {
		return err
	}
	c.Cookie(response.BuildAuthCookie(session))
	return c.Status(fiber.StatusCreated).JSON(response.Success("team joined successfully", session))
}

func New(teamUsecase teamDom.Usecase) TeamHandler {
	return &teamHandler{
		teamUc: teamUsecase,
	}
}
