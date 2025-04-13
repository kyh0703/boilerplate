package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/entity"
)

//counterfeiter:generate . EdgeRepository
type EdgeRepository interface {
	CreateOne(ctx context.Context, arg entity.Edge) (*entity.Edge, error)
	FindOne(ctx context.Context, id int64) (*entity.Edge, error)
	UpdateOne(ctx context.Context, arg entity.Edge) error
	DeleteOne(ctx context.Context, id int64) error
}
