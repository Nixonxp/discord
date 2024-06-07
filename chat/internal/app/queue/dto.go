package queue

type MessageKafkaMessage struct {
	Id      string `json:"id"`
	Text    string `json:"text"`
	ChatId  string `json:"chat_id"`
	OwnerId string `json:"owner_id"`
}
