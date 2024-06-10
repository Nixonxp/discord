package usecases

type CreateServerRequest struct {
	Name          string
	CurrentUserId string
}

type SearchServerRequest struct {
	Name string
}

type SubscribeServerRequest struct {
	ServerId      string
	CurrentUserId string
}

type UnsubscribeServerRequest struct {
	ServerId      string
	CurrentUserId string
}

type SearchServerByUserIdRequest struct {
	UserId string
}

type InviteUserToServerRequest struct {
	UserId   string
	ServerId string
}

type PublishMessageOnServerRequest struct {
	ServerId      string
	Text          string
	CurrentUserId string
}

type GetMessagesFromServerRequest struct {
	ServerId      string
	CurrentUserId string
}
