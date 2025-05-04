package oauth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/kyh0703/template/configs"
	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/service/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/kakao"

	dto "github.com/kyh0703/template/internal/core/dto/auth"
)

type oauthService struct {
	config          *configs.Config
	authService     auth.Service
	usersRepository repository.UsersRepository
	oauthRepository repository.OAuthRepository
	providers       map[Provider]*providerConfig
}

func NewOAuthService(
	config *configs.Config,
	authService auth.Service,
	usersRepository repository.UsersRepository,
	oauthRepository repository.OAuthRepository,
) Service {
	providers := map[Provider]*providerConfig{
		Google: {
			oauth2Config: &oauth2.Config{
				ClientID:     config.Auth.Google.ClientID,
				ClientSecret: config.Auth.Google.ClientSecret,
				RedirectURL:  config.Auth.Google.RedirectURL,
				Scopes:       config.Auth.Google.Scopes,
				Endpoint:     google.Endpoint,
			},
			userInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
			mapUserInfo: mapGoogleUserInfo,
		},
		Kakao: {
			oauth2Config: &oauth2.Config{
				ClientID:     config.Auth.Kakao.ClientID,
				ClientSecret: config.Auth.Kakao.ClientSecret,
				RedirectURL:  config.Auth.Kakao.RedirectURL,
				Scopes:       config.Auth.Kakao.Scopes,
				Endpoint:     kakao.Endpoint,
			},
			userInfoURL: "https://kapi.kakao.com/v2/user/me",
			mapUserInfo: mapKakaoUserInfo,
		},
		Github: {
			oauth2Config: &oauth2.Config{
				ClientID:     config.Auth.Github.ClientID,
				ClientSecret: config.Auth.Github.ClientSecret,
				RedirectURL:  config.Auth.Github.RedirectURL,
				Scopes:       config.Auth.Github.Scopes,
				Endpoint:     github.Endpoint,
			},
			userInfoURL: "https://api.github.com/user",
			mapUserInfo: mapGithubUserInfo,
		},
	}

	return &oauthService{
		config:          config,
		authService:     authService,
		usersRepository: usersRepository,
		oauthRepository: oauthRepository,
		providers:       providers,
	}
}

func mapGoogleUserInfo(data []byte) (*userInfo, error) {
	var info googleUserInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &userInfo{
		email:      info.Email,
		name:       info.Name,
		providerID: info.ID,
	}, nil
}

func mapKakaoUserInfo(data []byte) (*userInfo, error) {
	var info kakaoUserInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &userInfo{
		email:      info.KakaoAccount.Email,
		name:       info.Properties.Nickname,
		providerID: fmt.Sprintf("%d", info.ID),
	}, nil
}

func mapGithubUserInfo(data []byte) (*userInfo, error) {
	var info githubUserInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &userInfo{
		email:      info.Email,
		name:       info.Login,
		providerID: fmt.Sprintf("%d", info.ID),
	}, nil
}

func (s *oauthService) GenerateAuthURL(provider Provider, state string, redirectURL string) (string, error) {
	if _, err := s.oauthRepository.CreateState(context.Background(), model.CreateOAuthStateParams{
		State:       state,
		RedirectUrl: url.QueryEscape(redirectURL),
		ExpiresAt:   time.Now().Add(15 * time.Minute).Format(time.RFC3339),
	}); err != nil {
		return "", err
	}
	return s.providers[provider].oauth2Config.AuthCodeURL(state), nil
}

func (s *oauthService) GetRedirectURL(state string, token dto.Token) (string, error) {
	savedState, err := s.oauthRepository.FindByState(context.Background(), state)
	if err != nil {
		return "", err
	}

	decodedURL, err := url.QueryUnescape(savedState.RedirectUrl)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s?token=%s&expires_in=%d",
		decodedURL,
		token.Access.AccessToken,
		token.Access.AccessExpiresIn,
	), nil
}

func (s *oauthService) HandleCallback(ctx context.Context, provider Provider, code string, state string) (*dto.Token, error) {
	savedState, err := s.oauthRepository.FindByState(ctx, state)
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

	providerCfg := s.providers[provider]
	token, err := providerCfg.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := providerCfg.oauth2Config.Client(ctx, token)
	resp, err := client.Get(providerCfg.userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	userInfo, err := providerCfg.mapUserInfo(body)
	if err != nil {
		return nil, fmt.Errorf("failed mapping user info: %w", err)
	}

	user, err := s.usersRepository.FindByEmail(ctx, userInfo.email)
	if err != nil {
		user, err = s.usersRepository.CreateOne(ctx, model.CreateUserParams{
			Email:      userInfo.email,
			Name:       userInfo.name,
			Provider:   sql.NullString{String: string(provider), Valid: true},
			ProviderID: sql.NullString{String: userInfo.providerID, Valid: true},
			IsAdmin:    0,
		})
		if err != nil {
			return nil, fmt.Errorf("failed creating user: %w", err)
		}
	}

	authToken, err := s.authService.GenerateTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to validate and refresh token: %w", err)
	}

	return authToken, nil
}
