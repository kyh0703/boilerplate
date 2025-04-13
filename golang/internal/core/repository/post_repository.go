package repository

import (
	"context"
	"database/sql"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type PostRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewPostRepository(
	db *sql.DB,
	queries *model.Queries,
) repository.PostRepository {
	return &PostRepository{
		db:      db,
		queries: queries,
	}
}

func (p *PostRepository) CreateOne(ctx context.Context, arg model.CreatePostParams) (model.Post, error) {
	return p.queries.CreatePost(ctx, arg)
}

func (p *PostRepository) FindOne(ctx context.Context, id int64) (model.Post, error) {
	return p.queries.GetPost(ctx, id)
}

func (p *PostRepository) GetList(ctx context.Context, userID int64) ([]model.Post, error) {
	return p.queries.ListPosts(ctx, userID)
}

func (p *PostRepository) UpdateOne(ctx context.Context, arg model.PatchPostParams) error {
	return p.queries.PatchPost(ctx, arg)
}

func (p *PostRepository) DeleteOne(ctx context.Context, id int64) error {
	return p.queries.DeletePost(ctx, id)
}
