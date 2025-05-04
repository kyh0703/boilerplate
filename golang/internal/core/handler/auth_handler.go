package handler

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/internal/core/middleware"
	"github.com/kyh0703/template/internal/core/service/auth"
	"github.com/kyh0703/template/internal/pkg/response"

	dto "github.com/kyh0703/template/internal/core/dto/auth"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . AuthHandler
type AuthHandler interface {
	Handler
	Whoami(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type authHandler struct {
	validate       *validator.Validate
	authMiddleware middleware.AuthMiddleware
	authService    auth.Service
}

func NewAuthHandler(
	validate *validator.Validate,
	authMiddleware middleware.AuthMiddleware,
	authService auth.Service,
) AuthHandler {
	return &authHandler{
		validate:       validate,
		authMiddleware: authMiddleware,
		authService:    authService,
	}
}

func (a *authHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodGet, "/auth/whoami",
			a.authMiddleware.CurrentUser(), a.Whoami),
		Mapping(fiber.MethodPost, "/auth/register",
			a.Register),
		Mapping(fiber.MethodPost, "/auth/login",
			a.Login),
		Mapping(fiber.MethodPost, "/auth/logout",
			a.Logout),
		Mapping(fiber.MethodPost, "/auth/refresh",
			a.Refresh),
	}
}

func (a *authHandler) Whoami(c *fiber.Ctx) error {
	user := c.Locals("user")
	return response.Success(c, fiber.StatusOK, user)
}

func (a *authHandler) Register(c *fiber.Ctx) error {
	var register dto.RegisterDto
	if err := c.BodyParser(&register); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := a.validate.Struct(register); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := a.authService.Register(c.Context(), &register)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token.Refresh.RefreshToken,
		Expires:  time.Unix(token.Refresh.RefreshExpiresIn, 0),
		HTTPOnly: true,
		Secure:   false,
	})

	return response.Success(c, fiber.StatusOK, token.Access)
}

func (a *authHandler) Login(c *fiber.Ctx) error {
	var login dto.LoginDto
	if err := c.BodyParser(&login); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := a.validate.Struct(login); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := a.authService.Login(c.Context(), &login)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token.Refresh.RefreshToken,
		Expires:  time.Unix(token.Refresh.RefreshExpiresIn, 0),
		HTTPOnly: true,
		Secure:   false,
	})

	return response.Success(c, fiber.StatusOK, token.Access)
}

func (a *authHandler) Logout(c *fiber.Ctx) error {
	return response.Success(c, fiber.StatusOK, nil)
}

func (a *authHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("token")
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Refresh token not found")
	}

	token, err := a.authService.Refresh(c.Context(), refreshToken)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token.Refresh.RefreshToken,
		Expires:  time.Unix(token.Refresh.RefreshExpiresIn, 0),
		HTTPOnly: true,
		Secure:   false,
	})

	return response.Success(c, fiber.StatusOK, token.Access)
}
