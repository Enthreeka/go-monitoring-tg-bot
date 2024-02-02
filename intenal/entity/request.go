package entity

import (
	"fmt"
	"time"
)

type Request struct {
	ID                int       `json:"id"`
	UserID            int64     `json:"user_id"`
	ChannelTelegramID int64     `json:"channel_tg_id"`
	StatusRequest     string    `json:"status_request"`
	DateRequest       time.Time `json:"date_request"`
}

func (r Request) String() string {
	return fmt.Sprintf("(id: %d | user_id: %d | ChannelTelegramID: %d | status_request: %s)", r.ID, r.UserID,
		r.ChannelTelegramID, r.StatusRequest)
}
