package handler

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/dto/post"
	"github.com/kyh0703/template/internal/core/middleware"
	"github.com/kyh0703/template/internal/pkg/db"
	"github.com/kyh0703/template/internal/pkg/response"
)

//counterfeiter:generate . PostHandler
type PostHandler interface {
	Handler
	CreateOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	FindList(c *fiber.Ctx) error
}

type postHandler struct {
	validate       *validator.Validate
	authMiddleware middleware.AuthMiddleware
	postRepository repository.PostRepository
}

func NewPostHandler(
	validate *validator.Validate,
	authMiddleware middleware.AuthMiddleware,
	postRepository repository.PostRepository,
) PostHandler {
	return &postHandler{
		validate:       validate,
		authMiddleware: authMiddleware,
		postRepository: postRepository,
	}
}

func (h *postHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodPost, "/post",
			h.authMiddleware.CurrentUser(), h.CreateOne),
		Mapping(fiber.MethodGet, "/post/:id",
			h.authMiddleware.CurrentUser(), h.FindOne),
		Mapping(fiber.MethodPatch, "/post/:id",
			h.authMiddleware.CurrentUser(), h.UpdateOne),
		Mapping(fiber.MethodDelete, "/post/:id",
			h.authMiddleware.CurrentUser(), h.DeleteOne),
		Mapping(fiber.MethodPost, "/posts",
			h.authMiddleware.CurrentUser(), h.FindList),
	}
}

func (h *postHandler) CreateOne(c *fiber.Ctx) error {
	var req post.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var arg model.CreatePostParams
	copier.Copy(&arg, &req)

	user := c.Locals("user").(model.User)
	arg.UserID = user.ID

	newPost, err := h.postRepository.CreateOne(c.Context(), arg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var res post.PostResponse
	copier.Copy(&res, &newPost)

	return response.Success(c, fiber.StatusCreated, res)
}

func (h *postHandler) FindOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	findPost, err := h.postRepository.FindOne(c.Context(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var res post.PostResponse
	copier.Copy(&res, &findPost)

	return response.Success(c, fiber.StatusOK, res)
}

func (h *postHandler) UpdateOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req post.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := h.postRepository.FindOne(c.Context(), int64(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := h.postRepository.UpdateOne(c.Context(), model.PatchPostParams{
		ID:      int64(id),
		Title:   db.ToNullString(req.Name),
		Content: db.ToNullString(req.Description),
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *postHandler) DeleteOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.postRepository.DeleteOne(c.Context(), int64(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *postHandler) FindList(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)

	postList, err := h.postRepository.GetList(c.Context(), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var res []post.PostResponse
	copier.Copy(&res, &postList)

	return response.Success(c, fiber.StatusOK, res)
}
