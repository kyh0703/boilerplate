package repository

import (
	"context"
	"database/sql"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type postsRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewPostsRepository(
	db *sql.DB,
	queries *model.Queries,
) repository.PostsRepository {
	return &postsRepository{
		db:      db,
		queries: queries,
	}
}

func (p *postsRepository) CreateOne(ctx context.Context, arg model.CreatePostParams) (model.Post, error) {
	return p.queries.CreatePost(ctx, arg)
}

func (p *postsRepository) FindByID(ctx context.Context, id int64) (model.Post, error) {
	return p.queries.GetPost(ctx, id)
}

func (p *postsRepository) Pagination(ctx context.Context, userID int64) ([]model.Post, error) {
	return p.queries.ListPosts(ctx, userID)
}

func (p *postsRepository) UpdateOne(ctx context.Context, arg model.PatchPostParams) error {
	return p.queries.PatchPost(ctx, arg)
}

func (p *postsRepository) DeleteOne(ctx context.Context, id int64) error {
	return p.queries.DeletePost(ctx, id)
}
