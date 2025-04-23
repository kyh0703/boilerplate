package repository

import (
	"context"
	"database/sql"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type oauthStateRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewOAuthStateRepository(
	db *sql.DB,
	queries *model.Queries,
) repository.OAuthStateRepository {
	return &oauthStateRepository{
		db:      db,
		queries: queries,
	}
}

func (r *oauthStateRepository) CreateState(ctx context.Context, arg model.CreateOAuthStateParams) (model.OauthState, error) {
	return r.queries.CreateOAuthState(ctx, arg)
}

func (r *oauthStateRepository) GetState(ctx context.Context, state string) (model.OauthState, error) {
	return r.queries.GetOAuthState(ctx, state)
}

func (r *oauthStateRepository) DeleteState(ctx context.Context, state string) error {
	return r.queries.DeleteOAuthState(ctx, state)
}
