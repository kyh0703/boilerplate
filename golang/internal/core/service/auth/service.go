package auth

import (
	"context"

	"github.com/kyh0703/template/internal/core/dto/auth"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Service
type Service interface {
	Register(ctx context.Context, req *auth.Register) (*auth.Token, error)
	Login(ctx context.Context, req *auth.Login) (*auth.Token, error)
	Logout(ctx context.Context) error
	Refresh(ctx context.Context, refreshToken string) (*auth.Token, error)
}
