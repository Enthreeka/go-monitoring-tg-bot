package entity

import (
	"fmt"
	"time"
)

type User struct {
	ID          int64     `json:"id,omitempty"`
	UsernameTg  string    `json:"tg_username"`
	Phone       *string   `json:"phone,omitempty"`
	ChannelFrom *string   `json:"channel_from,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Role        string    `json:"user_role"`
	BlockedBot  bool      `json:"blocked_bot"`

	ChannelTelegramID int64 `json:"channel_telegram_id,omitempty"`
}

func (u User) String() string {
	return fmt.Sprintf("(id: %d | tg_username: %s | channel_from: %v | created_at: %v | role: %s)",
		u.ID, u.UsernameTg, u.ChannelFrom, u.CreatedAt, u.Role)
}
