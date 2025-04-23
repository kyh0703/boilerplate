package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
	"github.com/kyh0703/template/internal/core/dto/auth"
	"github.com/kyh0703/template/internal/pkg/jwt"
	"github.com/kyh0703/template/internal/pkg/password"
)

type authService struct {
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
}

func NewAuthService(
	authRepository repository.AuthRepository,
	userRepository repository.UserRepository,
) Service {
	return &authService{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (a *authService) generateNewTokens(ctx context.Context, user model.User) (*auth.Token, error) {
	accessExpire := time.Now().Add(jwt.AccessTokenExpireDuration)
	accessToken, err := jwt.GenerateToken(user.Email, accessExpire)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshExpire := time.Now().Add(jwt.RefreshTokenExpireDuration)
	refreshToken, err := jwt.GenerateToken(user.Email, refreshExpire)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if _, err := a.authRepository.CreateOne(ctx, model.CreateTokenParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresIn:    refreshExpire.Unix(),
	}); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &auth.Token{
		Access: auth.AccessToken{
			AccessToken:     accessToken,
			AccessExpiresIn: accessExpire.Unix(),
		},
		Refresh: auth.RefreshToken{
			RefreshToken:     refreshToken,
			RefreshExpiresIn: refreshExpire.Unix(),
		},
	}, nil
}

func (a *authService) Register(ctx context.Context, req *auth.Register) (*auth.Token, error) {
	_, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if req.Password != req.ConfirmPassword {
		return nil, fiber.NewError(fiber.StatusNotFound, "password and password confirm do not match")
	}

	hash, err := password.Hashed(req.Password)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	req.Password = hash

	var newUser model.CreateUserParams
	copier.Copy(&newUser, req)

	createdUser, err := a.userRepository.CreateOne(ctx, newUser)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return a.generateNewTokens(ctx, createdUser)
}

func (a *authService) Login(ctx context.Context, req *auth.Login) (*auth.Token, error) {
	user, err := a.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	ok, err := password.Compare(user.Password.String, req.Password)
	if err != nil || !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid email or password")
	}

	token, err := a.authRepository.FindOneByUserID(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return a.generateNewTokens(ctx, user)
	}

	expire := time.Unix(token.ExpiresIn, 0)
	if expire.After(time.Now()) {
		if err := a.authRepository.DeleteOne(ctx, token.ID); err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return a.generateNewTokens(ctx, user)
	}

	accessExpire := time.Now().Add(jwt.AccessTokenExpireDuration)
	accessToken, err := jwt.GenerateToken(req.Email, accessExpire)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &auth.Token{
		Access: auth.AccessToken{
			AccessToken:     accessToken,
			AccessExpiresIn: accessExpire.Unix(),
		},
		Refresh: auth.RefreshToken{
			RefreshToken:     token.RefreshToken,
			RefreshExpiresIn: token.ExpiresIn,
		},
	}, nil
}

func (a *authService) Logout(ctx context.Context) error {
	return nil
}

func (a *authService) Refresh(ctx context.Context, refreshToken string) (*auth.Token, error) {
	mapClaims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	email := mapClaims["email"].(string)
	if email == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := a.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	token, err := a.authRepository.FindOneByUserID(ctx, user.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	expire := time.Unix(token.ExpiresIn, 0)
	if expire.Before(time.Now()) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	accessExpire := time.Now().Add(jwt.AccessTokenExpireDuration)
	accessToken, err := jwt.GenerateToken(email, accessExpire)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &auth.Token{
		Access: auth.AccessToken{
			AccessToken:     accessToken,
			AccessExpiresIn: accessExpire.Unix(),
		},
		Refresh: auth.RefreshToken{
			RefreshToken:     token.RefreshToken,
			RefreshExpiresIn: token.ExpiresIn,
		},
	}, nil
}
