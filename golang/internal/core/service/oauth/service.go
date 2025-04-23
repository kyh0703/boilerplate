package oauth

import (
	"context"

	"github.com/kyh0703/template/internal/core/domain/model"
)

type Service interface {
	GetGoogleAuthURL(state string, redirectURL string) string
	HandleGoogleCallback(ctx context.Context, code string, state string) (*model.User, error)
	GetKakaoAuthURL(state string, redirectURL string) string
	HandleKakaoCallback(ctx context.Context, code string, state string) (*model.User, error)
}
