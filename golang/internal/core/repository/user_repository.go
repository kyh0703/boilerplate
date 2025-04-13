package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type userRepository struct {
	queries *model.Queries
}

func NewUserRepository(
	queries *model.Queries,
) repository.UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (u *userRepository) CreateOne(ctx context.Context, arg model.CreateUserParams) (model.User, error) {
	return u.queries.CreateUser(ctx, arg)
}

func (u *userRepository) FindOne(ctx context.Context, id int64) (model.User, error) {
	return u.queries.GetUser(ctx, id)
}

func (u *userRepository) FindOneByEmail(ctx context.Context, email string) (model.User, error) {
	return u.queries.GetUserByEmail(ctx, email)
}

func (u *userRepository) UpdateOne(ctx context.Context, arg model.UpdateUserParams) error {
	return u.queries.UpdateUser(ctx, arg)
}

func (u *userRepository) DeleteOne(ctx context.Context, id int64) error {
	return u.queries.DeleteUser(ctx, id)
}
