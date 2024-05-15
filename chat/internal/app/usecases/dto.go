package usecases

type SendUserPrivateMessageRequest struct {
	UserId uint64
	Text   string
}

type GetUserPrivateMessagesRequest struct {
	UserId uint64
}
