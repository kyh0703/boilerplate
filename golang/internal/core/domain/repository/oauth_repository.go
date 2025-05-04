package repository

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

type OAuthRepository interface {
	CreateState(ctx context.Context, arg model.CreateOAuthStateParams) (model.OauthState, error)
	FindByState(ctx context.Context, state string) (model.OauthState, error)
	DeleteState(ctx context.Context, state string) error
}
