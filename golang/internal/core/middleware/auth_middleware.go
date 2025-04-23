package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/pkg/jwt"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . AuthMiddleware
type AuthMiddleware interface {
	CurrentUser() fiber.Handler
	RequireAdmin() fiber.Handler
}

type authMiddleware struct {
	userRepository repository.UserRepository
}

func NewAuthMiddleware(
	userRepository repository.UserRepository,
) AuthMiddleware {
	return &authMiddleware{
		userRepository: userRepository,
	}
}

func (a *authMiddleware) CurrentUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken := c.Get("Authorization")
		if bearerToken == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		var tokenString string
		strArr := strings.Split(bearerToken, " ")
		if len(strArr) == 2 {
			tokenString = strArr[1]
		}

		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		mapClaims, err := jwt.ParseToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		email := mapClaims["email"].(string)
		if email == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		user, err := a.userRepository.FindByEmail(c.Context(), email)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		c.Locals("user", user)
		return c.Next()
	}
}

func (a *authMiddleware) RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(model.User)

		if user.IsAdmin == 0 {
			return fiber.NewError(fiber.StatusForbidden, "admin permission required")
		}

		return c.Next()
	}
}
