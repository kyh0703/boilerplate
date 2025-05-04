package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . AuthRepository
type AuthRepository interface {
	CreateOne(ctx context.Context, arg model.CreateTokenParams) (model.Token, error)
	FindByID(ctx context.Context, id int64) (model.Token, error)
	FindByUserID(ctx context.Context, userID int64) (model.Token, error)
	GetList(ctx context.Context) ([]model.Token, error)
	UpdateOne(ctx context.Context, arg model.UpdateTokenParams) error
	DeleteOne(ctx context.Context, id int64) error
}
