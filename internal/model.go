package internal

import "time"

// User model for persistence and auth
type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // bcrypt
}

// Chat message model
type Message struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to,omitempty"`
	Group     string    `json:"group,omitempty"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // direct | group | broadcast
	Timestamp time.Time `json:"timestamp"`
}
