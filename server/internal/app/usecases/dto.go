package usecases

type CreateServerRequest struct {
	Name string
}

type SearchServerRequest struct {
	Id   uint64
	Name string
}

type SubscribeServerRequest struct {
	ServerId uint64
}

type UnsubscribeServerRequest struct {
	ServerId uint64
}

type SearchServerByUserIdRequest struct {
	UserId uint64
}

type InviteUserToServerRequest struct {
	UserId   uint64
	ServerId uint64
}

type PublishMessageOnServerRequest struct {
	ServerId uint64
	Text     string
}

type GetMessagesFromServerRequest struct {
	ServerId uint64
}
