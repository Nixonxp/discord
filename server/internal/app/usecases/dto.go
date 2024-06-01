package usecases

type CreateServerRequest struct {
	Name string
}

type SearchServerRequest struct {
	Name string
}

type SubscribeServerRequest struct {
	ServerId string
}

type UnsubscribeServerRequest struct {
	ServerId string
}

type SearchServerByUserIdRequest struct {
	UserId string
}

type InviteUserToServerRequest struct {
	UserId   string
	ServerId string
}

type PublishMessageOnServerRequest struct {
	ServerId string
	Text     string
}

type GetMessagesFromServerRequest struct {
	ServerId string
}
