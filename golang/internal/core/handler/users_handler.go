package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/middleware"
)

//counterfeiter:generate . UsersHandler
type UsersHandler interface {
	Handler
	FindOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
}

type usersHandler struct {
	authMiddleware  middleware.AuthMiddleware
	usersRepository repository.UsersRepository
}

func NewUsersHandler(
	authMiddleware middleware.AuthMiddleware,
	usersRepository repository.UsersRepository,
) UsersHandler {
	return &usersHandler{
		authMiddleware:  authMiddleware,
		usersRepository: usersRepository,
	}
}

func (u *usersHandler) Table() []Mapper {
	return []Mapper{
		Mapping(
			fiber.MethodGet,
			"/user/:id",
			u.authMiddleware.CurrentUser(),
			u.FindOne,
		),
		Mapping(
			fiber.MethodPatch,
			"/user/:id",
			u.authMiddleware.CurrentUser(),
			u.UpdateOne,
		),
		Mapping(
			fiber.MethodDelete,
			"/user/:id",
			u.authMiddleware.CurrentUser(),
			u.DeleteOne,
		),
	}
}

func (u *usersHandler) FindOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (u *usersHandler) UpdateOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (u *usersHandler) DeleteOne(c *fiber.Ctx) error {
	panic("unimplemented")
}
