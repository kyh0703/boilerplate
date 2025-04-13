package usecase

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . UserUsecase
type UserUsecase interface {
	Create(ctx context.Context, email, username, password string)
	Get(ctx context.Context, email string)
	GetById(ctx context.Context, id int)
	Update(ctx context.Context, id int, param model.UpdateUserParams) error
	Remove(ctx context.Context, id int) error
}
