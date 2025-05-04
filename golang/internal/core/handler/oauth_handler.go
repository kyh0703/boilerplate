package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/service/oauth"
)

//counterfeiter:generate . OAuthHandler
type OAuthHandler interface {
	Handler
	GoogleLogin(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
	GithubLogin(c *fiber.Ctx) error
	GithubCallback(c *fiber.Ctx) error
	KakaoLogin(c *fiber.Ctx) error
	KakaoCallback(c *fiber.Ctx) error
}

type oauthHandler struct {
	oauthService    oauth.Service
	oauthRepository repository.OAuthRepository
}

func NewOAuthHandler(
	oauthService oauth.Service,
	oauthRepository repository.OAuthRepository,
) OAuthHandler {
	return &oauthHandler{
		oauthService:    oauthService,
		oauthRepository: oauthRepository,
	}
}

func (o *oauthHandler) Table() []Mapper {
	return []Mapper{
		Mapping(
			fiber.MethodGet,
			"/auth/google/login",
			o.GoogleLogin,
		),
		Mapping(
			fiber.MethodGet,
			"/auth/google/callback",
			o.GoogleCallback,
		),
		Mapping(
			fiber.MethodGet,
			"/auth/github/login",
			o.GithubLogin,
		),
		Mapping(
			fiber.MethodGet,
			"/auth/kakao/login",
			o.KakaoLogin,
		),
	}
}

func (o *oauthHandler) handleLogin(c *fiber.Ctx, provider oauth.Provider) error {
	state := uuid.New().String()
	redirectURL := c.Query("redirect_url", "/")

	authURL, err := o.oauthService.GenerateAuthURL(provider, state, redirectURL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Redirect(authURL)
}

func (o *oauthHandler) handleCallback(c *fiber.Ctx, provider oauth.Provider) error {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing code or state")
	}

	token, err := o.oauthService.HandleCallback(c.Context(), provider, code, state)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	redirectURL, err := o.oauthService.GetRedirectURL(state, *token)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token.Refresh.RefreshToken,
		Expires:  time.Unix(token.Refresh.RefreshExpiresIn, 0),
		HTTPOnly: true,
		Secure:   false,
	})

	return c.Redirect(redirectURL)
}

func (o *oauthHandler) GoogleLogin(c *fiber.Ctx) error {
	return o.handleLogin(c, oauth.Google)
}

func (o *oauthHandler) GoogleCallback(c *fiber.Ctx) error {
	return o.handleCallback(c, oauth.Google)
}

func (o *oauthHandler) GithubLogin(c *fiber.Ctx) error {
	return o.handleLogin(c, oauth.Github)
}

func (o *oauthHandler) GithubCallback(c *fiber.Ctx) error {
	return o.handleCallback(c, oauth.Github)
}

func (o *oauthHandler) KakaoLogin(c *fiber.Ctx) error {
	return o.handleLogin(c, oauth.Kakao)
}

func (o *oauthHandler) KakaoCallback(c *fiber.Ctx) error {
	return o.handleCallback(c, oauth.Kakao)
}
