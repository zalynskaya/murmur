package entity

import "time"

type Message struct {
	MessageID int       `json:"message_id,omitempty"` // MessageID = Id
	ChatID    int       `json:"chat_id,omitempty"`
	AuthorID  string    `json:"author_id,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
