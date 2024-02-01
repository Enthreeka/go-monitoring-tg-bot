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
	return fmt.Sprintf("(id: %d | tg_username: %s | phone: %s | channel_from: %s | created_at: %v | role: %s)",
		u.ID, u.UsernameTg, *u.Phone, *u.ChannelFrom, u.CreatedAt, u.Role)
}
