package domain

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Score    uint   `json:"score,omitempty"`
	KahootId uint	`json:"kahoot_id"`
}

func NewUser(username string) *User {
	return &User{Username: username}
}
