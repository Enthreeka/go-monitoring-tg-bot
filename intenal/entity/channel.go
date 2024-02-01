package entity

import "fmt"

type Channel struct {
	ID          int     `json:"id"`
	TelegramID  int64   `json:"tg_id"`
	ChannelName string  `json:"channel_name"`
	ChannelURL  *string `json:"channel_url"`

	Status string `json:"status"`
}

func (c Channel) String() string {
	return fmt.Sprintf("(id: %d | tg_id: %d | channel_name: %s | ChannelURL: %s | Status: %s)", c.ID,
		c.TelegramID, c.ChannelName, *c.ChannelURL, c.Status)
}
