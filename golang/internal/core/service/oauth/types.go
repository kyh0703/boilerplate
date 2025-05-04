package oauth

import "golang.org/x/oauth2"

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type kakaoUserInfo struct {
	ID         int64 `json:"id"`
	Properties struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		Email string `json:"email"`
	} `json:"kakao_account"`
}

type githubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}

type providerConfig struct {
	oauth2Config *oauth2.Config
	userInfoURL  string
	mapUserInfo  func([]byte) (*userInfo, error)
}

type userInfo struct {
	email      string
	name       string
	providerID string
}
