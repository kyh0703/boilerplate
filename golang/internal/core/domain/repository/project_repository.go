package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

//counterfeiter:generate . ProjectRepository
type ProjectRepository interface {
	CreateOne(ctx context.Context, arg model.CreateProjectParams) (model.Project, error)
	FindOne(ctx context.Context, id int64) (model.Project, error)
	GetList(ctx context.Context, userID int64) ([]model.Project, error)
	UpdateOne(ctx context.Context, arg model.PatchProjectParams) error
	DeleteOne(ctx context.Context, id int64) error
}
