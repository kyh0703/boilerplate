package auth

type AccessToken struct {
	AccessToken     string `json:"accessToken"`
	AccessExpiresIn int64  `json:"accessExpiresIn"`
}

type RefreshToken struct {
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}

type Token struct {
	Access  AccessToken
	Refresh RefreshToken
}
