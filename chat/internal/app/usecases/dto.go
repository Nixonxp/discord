package usecases

type SendUserPrivateMessageRequest struct {
	UserId      string
	Text        string
	CurrentUser string
}

type GetUserPrivateMessagesRequest struct {
	UserId      string
	CurrentUser string
}

type MessageDto struct {
	Id      string `json:"id"`
	Text    string `json:"text"`
	ChatId  string `json:"chat_id"`
	OwnerId string `json:"owner_id"`
}

type CreatePrivateChatRequest struct {
	UserId      string
	CurrentUser string
}

type SendServerMessageRequest struct {
	ServerId string
	Text     string
}
type GetServerMessageRequest struct {
	ServerId string
}
