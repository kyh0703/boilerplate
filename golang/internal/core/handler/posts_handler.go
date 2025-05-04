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

//counterfeiter:generate . PostsHandler
type PostsHandler interface {
	Handler
	CreateOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
}

type postsHandler struct {
	validate       *validator.Validate
	authMiddleware middleware.AuthMiddleware
	postRepository repository.PostsRepository
}

func NewPostsHandler(
	validate *validator.Validate,
	authMiddleware middleware.AuthMiddleware,
	postRepository repository.PostsRepository,
) PostsHandler {
	return &postsHandler{
		validate:       validate,
		authMiddleware: authMiddleware,
		postRepository: postRepository,
	}
}

func (p *postsHandler) Table() []Mapper {
	return []Mapper{
		Mapping(fiber.MethodPost, "/post",
			p.authMiddleware.CurrentUser(), p.CreateOne),
		Mapping(fiber.MethodGet, "/post/:id",
			p.authMiddleware.CurrentUser(), p.FindOne),
		Mapping(fiber.MethodPatch, "/post/:id",
			p.authMiddleware.CurrentUser(), p.UpdateOne),
		Mapping(fiber.MethodDelete, "/post/:id",
			p.authMiddleware.CurrentUser(), p.DeleteOne),
		Mapping(fiber.MethodGet, "/posts",
			p.authMiddleware.CurrentUser(), p.FindAll),
	}
}

func (p *postsHandler) CreateOne(c *fiber.Ctx) error {
	var req post.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := p.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var arg model.CreatePostParams
	copier.Copy(&arg, &req)

	user := c.Locals("user").(model.User)
	arg.UserID = user.ID

	newPost, err := p.postRepository.CreateOne(c.Context(), arg)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var res post.PostDto
	copier.Copy(&res, &newPost)

	return response.Success(c, fiber.StatusCreated, res)
}

func (p *postsHandler) FindOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	findPost, err := p.postRepository.FindByID(c.Context(), int64(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var res post.PostDto
	copier.Copy(&res, &findPost)

	return response.Success(c, fiber.StatusOK, res)
}

func (p *postsHandler) UpdateOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req post.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := p.validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err := p.postRepository.FindByID(c.Context(), int64(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := p.postRepository.UpdateOne(c.Context(), model.PatchPostParams{
		ID:      int64(id),
		Title:   db.ToNullString(req.Name),
		Content: db.ToNullString(req.Description),
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (p *postsHandler) DeleteOne(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := p.postRepository.DeleteOne(c.Context(), int64(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (p *postsHandler) FindAll(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)

	postList, err := p.postRepository.Pagination(c.Context(), user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var res []post.PostDto
	copier.Copy(&res, &postList)

	return response.Success(c, fiber.StatusOK, res)
}
