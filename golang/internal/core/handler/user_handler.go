package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/middleware"
)

//counterfeiter:generate . UserHandler
type UserHandler interface {
	Handler
	FindOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
}

type userHandler struct {
	authMiddleware middleware.AuthMiddleware
	userRepository repository.UserRepository
}

func NewUserHandler(
	authMiddleware middleware.AuthMiddleware,
	userRepository repository.UserRepository,
) UserHandler {
	return &userHandler{
		authMiddleware: authMiddleware,
		userRepository: userRepository,
	}
}

func (u *userHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodGet, "/user/:id",
			u.authMiddleware.CurrentUser(), u.FindOne),
		Mapping(fiber.MethodPatch, "/user/:id",
			u.authMiddleware.CurrentUser(), u.UpdateOne),
		Mapping(fiber.MethodDelete, "/user/:id",
			u.authMiddleware.CurrentUser(), u.DeleteOne),
	}
}

func (u *userHandler) FindOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (u *userHandler) UpdateOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (u *userHandler) DeleteOne(c *fiber.Ctx) error {
	panic("unimplemented")
}
