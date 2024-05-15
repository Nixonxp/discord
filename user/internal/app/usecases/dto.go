package usecases

type UpdateUserRequest struct {
	Id       uint64
	Login    string
	Name     string
	Email    string
	Password string
}

type GetUserByLoginRequest struct {
	Login string
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
