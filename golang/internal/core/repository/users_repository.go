package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type usersRepository struct {
	queries *model.Queries
}

func NewUsersRepository(
	queries *model.Queries,
) repository.UsersRepository {
	return &usersRepository{
		queries: queries,
	}
}

func (u *usersRepository) CreateOne(ctx context.Context, arg model.CreateUserParams) (model.User, error) {
	return u.queries.CreateUser(ctx, arg)
}

func (u *usersRepository) FindByID(ctx context.Context, id int64) (model.User, error) {
	return u.queries.GetUser(ctx, id)
}

func (u *usersRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	return u.queries.GetUserByEmail(ctx, email)
}

func (u *usersRepository) UpdateOne(ctx context.Context, arg model.UpdateUserParams) error {
	return u.queries.UpdateUser(ctx, arg)
}

func (u *usersRepository) DeleteOne(ctx context.Context, id int64) error {
	return u.queries.DeleteUser(ctx, id)
}
