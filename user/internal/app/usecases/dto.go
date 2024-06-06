package usecases

type CreateUserRequest struct {
	Login    string
	Name     string
	Email    string
	Password string
}

type CreateOrGetUserRequest struct {
	Login          string
	Name           string
	Email          string
	Password       string
	OauthId        string
	AvatarPhotoUrl string
}

type UpdateUserRequest struct {
	Id             string
	Login          string
	Name           string
	Email          string
	AvatarPhotoUrl string
	CurrentUserId  string
}

type GetUserByLoginAndPasswordRequest struct {
	Login    string
	Password string
}

type GetUserByLoginRequest struct {
	Login string
}

type GetUserFriendsRequest struct {
	UserId string
}

type GetUserInvitesRequest struct {
	UserId string
}

type AddToFriendByUserIdRequest struct {
	UserId  string
	OwnerId string
}

type AcceptFriendInviteRequest struct {
	InviteId      string
	CurrentUserId string
}

type DeclineFriendInviteRequest struct {
	InviteId      string
	CurrentUserId string
}

type DeleteFromFriendRequest struct {
	FriendId      string
	CurrentUserId string
}
