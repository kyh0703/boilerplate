package auth

type Refresh struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
