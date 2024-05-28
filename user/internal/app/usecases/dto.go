package usecases

type CreateUserRequest struct {
	Login    string
	Name     string
	Email    string
	Password string
}

type UpdateUserRequest struct {
	Id             string
	Login          string
	Name           string
	Email          string
	AvatarPhotoUrl string
}

type GetUserByLoginAndPasswordRequest struct {
	Login    string
	Password string
}

type GetUserByLoginRequest struct {
	Login string
}

type GetUserFriendsRequest struct {
}

type GetUserInvitesRequest struct {
	UserId string
}

type AddToFriendByUserIdRequest struct {
	UserId string
}

type AcceptFriendInviteRequest struct {
	InviteId string
}

type DeclineFriendInviteRequest struct {
	InviteId string
}
