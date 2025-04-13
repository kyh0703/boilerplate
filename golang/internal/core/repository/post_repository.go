package repository

import (
	"context"
	"database/sql"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type projectRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewProjectRepository(
	db *sql.DB,
	queries *model.Queries,
) repository.ProjectRepository {
	return &projectRepository{
		db:      db,
		queries: queries,
	}
}

func (p *projectRepository) CreateOne(ctx context.Context, arg model.CreateProjectParams) (model.Project, error) {
	return p.queries.CreateProject(ctx, arg)
}

func (p *projectRepository) FindOne(ctx context.Context, id int64) (model.Project, error) {
	return p.queries.GetProject(ctx, id)
}

func (p *projectRepository) GetList(ctx context.Context, userID int64) ([]model.Project, error) {
	return p.queries.ListProjects(ctx, userID)
}

func (p *projectRepository) UpdateOne(ctx context.Context, arg model.PatchProjectParams) error {
	return p.queries.PatchProject(ctx, arg)
}

func (p *projectRepository) DeleteOne(ctx context.Context, id int64) error {
	return p.queries.DeleteProject(ctx, id)
}
