package usecases

type SendUserPrivateMessageRequest struct {
	UserId string
	Text   string
}

type GetUserPrivateMessagesRequest struct {
	UserId string
}
