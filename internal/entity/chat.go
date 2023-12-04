package entity

import "time"

type Chat struct {
	ChatID    int       `json:"id,omitempty"` // ChatID = ID
	Name      string    `json:"name,omitempty"`
	Users     []User    `json:"users,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
