package repository

import (
	"context"
	"database/sql"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type oauthRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewOAuthRepository(
	db *sql.DB,
	queries *model.Queries,
) repository.OAuthRepository {
	return &oauthRepository{
		db:      db,
		queries: queries,
	}
}

func (r *oauthRepository) CreateState(ctx context.Context, arg model.CreateOAuthStateParams) (model.OauthState, error) {
	return r.queries.CreateOAuthState(ctx, arg)
}

func (r *oauthRepository) FindByState(ctx context.Context, state string) (model.OauthState, error) {
	return r.queries.GetOAuthState(ctx, state)
}

func (r *oauthRepository) DeleteState(ctx context.Context, state string) error {
	return r.queries.DeleteOAuthState(ctx, state)
}
