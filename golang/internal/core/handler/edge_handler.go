package handler

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/middleware"
)

//counterfeiter:generate . EdgeHandler
type EdgeHandler interface {
	Handler
	CreateOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
}

type edgeHandler struct {
	validate       *validator.Validate
	authMiddleware middleware.AuthMiddleware
	edgeRepository repository.EdgeRepository
}

func NewEdgeHandler(
	validate *validator.Validate,
	authMiddleware middleware.AuthMiddleware,
	edgeRepository repository.EdgeRepository,
) EdgeHandler {
	return &edgeHandler{
		validate:       validate,
		authMiddleware: authMiddleware,
		edgeRepository: edgeRepository,
	}
}

func (h *edgeHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodPost, "/edge",
			h.authMiddleware.CurrentUser(), h.CreateOne),
		Mapping(fiber.MethodGet, "/edge/:id",
			h.authMiddleware.CurrentUser(), h.FindOne),
		Mapping(fiber.MethodPatch, "/edge/:id",
			h.authMiddleware.CurrentUser(), h.UpdateOne),
		Mapping(fiber.MethodDelete, "/edge/:id",
			h.authMiddleware.CurrentUser(), h.DeleteOne),
	}
}

func (h *edgeHandler) CreateOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *edgeHandler) FindOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *edgeHandler) DeleteOne(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *edgeHandler) UpdateOne(c *fiber.Ctx) error {
	panic("unimplemented")
}
