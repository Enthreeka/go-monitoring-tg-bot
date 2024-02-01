package entity

import (
	"fmt"
	"time"
)

type User struct {
	ID          int64     `json:"id"`
	UsernameTg  string    `json:"tg_username"`
	Phone       *string   `json:"phone"`
	ChannelFrom *string   `json:"channel_from"`
	CreatedAt   time.Time `json:"created_at"`
	Role        string    `json:"user_role"`
}

func (u User) String() string {
	return fmt.Sprintf("(id: %d | tg_username: %s | channel_from: %v | created_at: %v | role: %s)",
		u.ID, u.UsernameTg, u.ChannelFrom, u.CreatedAt, u.Role)
}
