package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

//counterfeiter:generate . PostRepository
type PostRepository interface {
	CreateOne(ctx context.Context, arg model.CreatePostParams) (model.Post, error)
	FindOne(ctx context.Context, id int64) (model.Post, error)
	GetList(ctx context.Context, userID int64) ([]model.Post, error)
	UpdateOne(ctx context.Context, arg model.PatchPostParams) error
	DeleteOne(ctx context.Context, id int64) error
}
