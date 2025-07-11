package models

type MessageRequest struct {
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Text             MessageText `json:"text"`
}

type MessageText struct {
	Body string `json:"body"`
}
