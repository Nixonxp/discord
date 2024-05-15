package models

type OauthLoginRequest struct {
}

type OauthLoginResult struct {
	Code string
}

type OauthLoginCallbackResult struct {
	User
}
