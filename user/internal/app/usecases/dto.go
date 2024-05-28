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

type GetUserFriendsRequest struct {
}

type AddToFriendByUserIdRequest struct {
	UserId uint64
}

type AcceptFriendInviteRequest struct {
	InviteId uint64
}

type DeclineFriendInviteRequest struct {
	InviteId uint64
}
