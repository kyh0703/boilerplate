package auth

import (
	"context"
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
	authRepository  repository.AuthRepository
	usersRepository repository.UsersRepository
}

func NewAuthService(
	authRepository repository.AuthRepository,
	usersRepository repository.UsersRepository,
) Service {
	return &authService{
		authRepository:  authRepository,
		usersRepository: usersRepository,
	}
}

func (a *authService) GenerateTokens(ctx context.Context, user model.User) (*auth.Token, error) {
	var (
		accessExpire  time.Time = time.Now().Add(jwt.AccessTokenExpireDuration)
		refreshExpire time.Time = time.Now().Add(jwt.RefreshTokenExpireDuration)
	)

	token, err := a.authRepository.FindByUserID(ctx, user.ID)
	if err == nil {
		expire := time.Unix(token.ExpiresIn, 0)
		if expire.After(time.Now()) {
			accessToken, err := jwt.GenerateToken(user.Email, accessExpire)
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

		if err := a.authRepository.DeleteOne(ctx, token.ID); err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	accessToken, err := jwt.GenerateToken(user.Email, accessExpire)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

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

func (a *authService) Register(ctx context.Context, req *auth.RegisterDto) (*auth.Token, error) {
	if _, err := a.usersRepository.FindByEmail(ctx, req.Email); err == nil {
		return nil, err
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("password and password confirm do not match")
	}

	hashedPassword, err := password.Hashed(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = hashedPassword

	var newUser model.CreateUserParams
	copier.Copy(&newUser, req)

	createdUser, err := a.usersRepository.CreateOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return a.GenerateTokens(ctx, createdUser)
}

func (a *authService) Login(ctx context.Context, req *auth.LoginDto) (*auth.Token, error) {
	user, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	ok, err := password.Compare(user.Password.String, req.Password)
	if err != nil || !ok {
		return nil, errors.New("invalid email or password")
	}

	return a.GenerateTokens(ctx, user)
}

func (a *authService) Logout(ctx context.Context) error {
	return nil
}

func (a *authService) Refresh(ctx context.Context, refreshToken string) (*auth.Token, error) {
	mapClaims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	email := mapClaims["email"].(string)
	if email == "" {
		return nil, errors.New("unauthorized")
	}

	user, err := a.usersRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return a.GenerateTokens(ctx, user)
}
