package oauth

import (
	"context"

	"github.com/kyh0703/template/internal/core/dto/auth"
)

type Provider string

const (
	Google Provider = "google"
	Kakao  Provider = "kakao"
	Github Provider = "github"
)

type Service interface {
	GenerateAuthURL(provider Provider, state string, redirectURL string) (string, error)
	GetRedirectURL(state string, token auth.Token) (string, error)
	HandleCallback(ctx context.Context, provider Provider, code string, state string) (*auth.Token, error)
}
