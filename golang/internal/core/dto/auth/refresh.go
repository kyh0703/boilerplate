package auth

type RefreshDto struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
