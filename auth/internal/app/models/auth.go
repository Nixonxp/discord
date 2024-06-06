package models

type OauthLoginRequest struct {
}

type OauthLoginResult struct {
	Code string
}

type LoginResult struct {
	Token        string
	RefreshToken string
}
