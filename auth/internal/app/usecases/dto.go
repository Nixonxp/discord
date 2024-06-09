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

type RefreshInfo struct {
	Token string
}

type OauthLoginRequest struct {
}

type OauthLoginCallbackRequest struct {
	Code  string
	State string
}

type GetOrCreateUserRequest struct {
	Login          string
	Name           string
	Email          string
	Password       string
	AvatarPhotoUrl string
	OauthId        string
}

type UserInfo struct {
	OauthId        string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvatarPhotoUrl string `json:"picture"`
}
