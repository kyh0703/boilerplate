package handler

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/dto/flows"
	"github.com/kyh0703/template/internal/core/middleware"
	"github.com/kyh0703/template/internal/pkg/db"
	"github.com/kyh0703/template/internal/pkg/response"
)

//counterfeiter:generate . flowHandler
type FlowHandler interface {
	Table() []Mapper
	CreateOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	Undo(c *fiber.Ctx) error
	Redo(c *fiber.Ctx) error
}

type flowHandler struct {
	validate       *validator.Validate
	AuthMiddleware middleware.AuthMiddleware
	flowRepository repository.FlowRepository
}

func NewFlowHandler(
	validate *validator.Validate,
	AuthMiddleware middleware.AuthMiddleware,
	flowRepository repository.FlowRepository,
) FlowHandler {
	return &flowHandler{
		validate:       validate,
		AuthMiddleware: AuthMiddleware,
		flowRepository: flowRepository,
	}
}

func (f *flowHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodPost, "/flow",
			f.AuthMiddleware.CurrentUser(), f.CreateOne),
		Mapping(fiber.MethodGet, "/flow/:id",
			f.AuthMiddleware.CurrentUser(), f.FindOne),
		Mapping(fiber.MethodPatch, "/flow/:id",
			f.AuthMiddleware.CurrentUser(), f.UpdateOne),
		Mapping(fiber.MethodDelete, "/flow/:id",
			f.AuthMiddleware.CurrentUser(), f.DeleteOne),
		Mapping(fiber.MethodPost, "/flow/:id/undo",
			f.AuthMiddleware.CurrentUser(), f.Undo),
		Mapping(fiber.MethodPost, "/flow/:id/redo",
			f.AuthMiddleware.CurrentUser(), f.Redo),
	}
}

func (f *flowHandler) CreateOne(c *fiber.Ctx) error {
	var req flows.CreateFlowRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := f.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var params model.CreateFlowParams
	copier.Copy(&params, &req)

	newFlow, err := f.flowRepository.CreateOne(c.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var res flows.FlowResponse
	copier.Copy(&res, &newFlow)

	return response.Success(c, fiber.StatusCreated, res)
}

func (f *flowHandler) FindOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	flow, err := f.flowRepository.FindOne(c.Context(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var res flows.FlowResponse
	copier.Copy(&res, &flow)

	return response.Success(c, fiber.StatusOK, res)
}

func (f *flowHandler) UpdateOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req flows.UpdateFlowRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := f.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := f.flowRepository.FindOne(c.Context(), int64(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := f.flowRepository.UpdateOne(c.Context(), model.PatchFlowParams{
		ID:          int64(id),
		Name:        db.ToNullString(req.Name),
		Description: db.ToNullString(req.Description),
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (f *flowHandler) DeleteOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = f.flowRepository.DeleteOne(c.Context(), int64(id)); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (f *flowHandler) Undo(c *fiber.Ctx) error {
	return nil
}

func (f *flowHandler) Redo(c *fiber.Ctx) error {
	return nil
}
