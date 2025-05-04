package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

//counterfeiter:generate . PostsRepository
type PostsRepository interface {
	CreateOne(ctx context.Context, arg model.CreatePostParams) (model.Post, error)
	FindByID(ctx context.Context, id int64) (model.Post, error)
	Pagination(ctx context.Context, userID int64) ([]model.Post, error)
	UpdateOne(ctx context.Context, arg model.PatchPostParams) error
	DeleteOne(ctx context.Context, id int64) error
}
