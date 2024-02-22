package entity

import "time"

type User struct {
	UserID    string    `json:"user_id,omitempty"`
	Username  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
