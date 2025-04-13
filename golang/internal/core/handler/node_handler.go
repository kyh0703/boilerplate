package handler

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/internal/core/middleware"
)

//counterfeiter:generate . NodeHandler
type NodeHandler interface {
	Handler
	CreateOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
}

type nodeHandler struct {
	validate       *validator.Validate
	authMiddleware middleware.AuthMiddleware
}

func NewNodeHandler(
	validate *validator.Validate,
	authMiddleware middleware.AuthMiddleware,
) NodeHandler {
	return &nodeHandler{
		validate:       validate,
		authMiddleware: authMiddleware,
	}
}

func (h *nodeHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodPost, "/node",
			h.authMiddleware.CurrentUser(), h.CreateOne),
		Mapping(fiber.MethodGet, "/node/:id",
			h.authMiddleware.CurrentUser(), h.FindOne),
		Mapping(fiber.MethodPatch, "/node/:id",
			h.authMiddleware.CurrentUser(), h.UpdateOne),
		Mapping(fiber.MethodDelete, "/node/:id",
			h.authMiddleware.CurrentUser(), h.DeleteOne),
	}
}

func (h *nodeHandler) CreateOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *nodeHandler) FindOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *nodeHandler) DeleteOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *nodeHandler) UpdateOne(c *fiber.Ctx) error {
	panic("unimplemented")
}
