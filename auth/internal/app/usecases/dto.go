package usecases

type RegisterUserInfo struct {
	Login    string
	Name     string
	Email    string
	Password string
}

type LoginUserInfo struct {
	Login    string
	Password string
}

type OauthLoginRequest struct {
}

type OauthLoginCallbackRequest struct {
	Code string
}
