package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

//counterfeiter:generate . UserRepository
type UserRepository interface {
	CreateOne(ctx context.Context, arg model.CreateUserParams) (model.User, error)
	FindOne(ctx context.Context, id int64) (model.User, error)
	FindOneByEmail(ctx context.Context, email string) (model.User, error)
	UpdateOne(ctx context.Context, arg model.UpdateUserParams) error
	DeleteOne(ctx context.Context, id int64) error
}
