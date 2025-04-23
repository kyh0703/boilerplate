package oauth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kyh0703/template/configs"
	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauthService struct {
	config          *configs.Config
	userRepository  repository.UserRepository
	stateRepository repository.OAuthStateRepository
	googleConf      *oauth2.Config
	kakaoConf       *oauth2.Config
}

func NewOAuthService(
	config *configs.Config,
	userRepository repository.UserRepository,
	stateRepository repository.OAuthStateRepository,
) Service {
	googleConf := &oauth2.Config{
		ClientID:     config.Auth.Google.ClientID,
		ClientSecret: config.Auth.Google.ClientSecret,
		RedirectURL:  config.Auth.Google.RedirectURL,
		Scopes:       config.Auth.Google.Scopes,
		Endpoint:     google.Endpoint,
	}

	kakaoConf := &oauth2.Config{
		ClientID:     config.Auth.Kakao.ClientID,
		ClientSecret: config.Auth.Kakao.ClientSecret,
		RedirectURL:  config.Auth.Kakao.RedirectURL,
		Scopes:       config.Auth.Kakao.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://kauth.kakao.com/oauth/authorize",
			TokenURL: "https://kauth.kakao.com/oauth/token",
		},
	}

	return &oauthService{
		config:          config,
		userRepository:  userRepository,
		stateRepository: stateRepository,
		googleConf:      googleConf,
		kakaoConf:       kakaoConf,
	}
}

func (s *oauthService) GetGoogleAuthURL(state string, redirectURL string) string {
	s.stateRepository.CreateState(context.Background(), model.CreateOAuthStateParams{
		State:       state,
		RedirectUrl: redirectURL,
		ExpiresAt:   time.Now().Add(15 * time.Minute).Format(time.RFC3339),
	})

	return s.googleConf.AuthCodeURL(state)
}

func (s *oauthService) HandleGoogleCallback(ctx context.Context, code string, state string) (*model.User, error) {
	savedState, err := s.stateRepository.GetState(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("invalid state: %w", err)
	}

	expiresAt, err := time.Parse(time.RFC3339, savedState.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time: %w", err)
	}
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("state expired")
	}

	token, err := s.googleConf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := s.googleConf.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed decoding user info: %w", err)
	}

	user, err := s.userRepository.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		user, err = s.userRepository.CreateOne(ctx, model.CreateUserParams{
			Email:      userInfo.Email,
			Name:       userInfo.Name,
			Provider:   sql.NullString{String: "google", Valid: true},
			ProviderID: sql.NullString{String: userInfo.ID, Valid: true},
			IsAdmin:    0,
		})
		if err != nil {
			return nil, fmt.Errorf("failed creating user: %w", err)
		}
	}

	return &user, nil
}

func (s *oauthService) GetKakaoAuthURL(state string, redirectURL string) string {
	s.stateRepository.CreateState(context.Background(), model.CreateOAuthStateParams{
		State:       state,
		RedirectUrl: redirectURL,
		ExpiresAt:   time.Now().Add(15 * time.Minute).Format(time.RFC3339),
	})

	return s.kakaoConf.AuthCodeURL(state)
}

func (s *oauthService) HandleKakaoCallback(ctx context.Context, code string, state string) (*model.User, error) {
	savedState, err := s.stateRepository.GetState(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("invalid state: %w", err)
	}

	expiresAt, err := time.Parse(time.RFC3339, savedState.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time: %w", err)
	}
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("state expired")
	}

	token, err := s.kakaoConf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := s.kakaoConf.Client(ctx, token)
	resp, err := client.Get("https://kapi.kakao.com/v2/user/me")
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo KakaoUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed decoding user info: %w", err)
	}

	user, err := s.userRepository.FindByEmail(ctx, userInfo.KakaoAccount.Email)
	if err != nil {
		user, err = s.userRepository.CreateOne(ctx, model.CreateUserParams{
			Email:      userInfo.KakaoAccount.Email,
			Name:       userInfo.Properties.Nickname,
			Provider:   sql.NullString{String: "kakao", Valid: true},
			ProviderID: sql.NullString{String: fmt.Sprintf("%d", userInfo.ID), Valid: true},
			IsAdmin:    0,
		})
		if err != nil {
			return nil, fmt.Errorf("failed creating user: %w", err)
		}
	}

	return &user, nil
}
