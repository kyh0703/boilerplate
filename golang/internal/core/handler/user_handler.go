package handler

import "github.com/gofiber/fiber/v2"

//counterfeiter:generate . UserHandler
type UserHandler interface {
	Handler
	CreateOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
}

type userHandler struct{}

func NewUserHandler() UserHandler {
	return &userHandler{}
}

func (u *userHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodPost, "/user", u.CreateOne),
		Mapping(fiber.MethodGet, "/user/:id", u.FindOne),
		Mapping(fiber.MethodPatch, "/user/:id", u.UpdateOne),
		Mapping(fiber.MethodDelete, "/user/:id", u.DeleteOne),
	}
}

func (u *userHandler) CreateOne(c *fiber.Ctx) error {
	panic("unimplemented")
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
