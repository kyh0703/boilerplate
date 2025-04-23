package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kyh0703/template/internal/core/service/oauth"
	"github.com/kyh0703/template/internal/pkg/response"
)

type OAuthHandler interface {
	Handler
	GoogleLogin(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
}

type oauthHandler struct {
	oauthService oauth.Service
}

func NewOAuthHandler(
	oauthService oauth.Service,
) OAuthHandler {
	return &oauthHandler{
		oauthService: oauthService,
	}
}

func (h *oauthHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodGet, "/auth/google/login", h.GoogleLogin),
		Mapping(fiber.MethodGet, "/auth/google/callback", h.GoogleCallback),
	}
}

func (h *oauthHandler) GoogleLogin(c *fiber.Ctx) error {
	state := uuid.New().String()
	redirectURL := c.Query("redirect_url", "/")

	authURL := h.oauthService.GetGoogleAuthURL(state, redirectURL)
	return c.Redirect(authURL)
}

func (h *oauthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing code or state")
	}

	user, err := h.oauthService.HandleGoogleCallback(c.Context(), code, state)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// TODO: Generate JWT token and set cookie

	return response.Success(c, fiber.StatusOK, user)
}
