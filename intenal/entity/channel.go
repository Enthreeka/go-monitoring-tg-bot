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
	var url string
	if c.ChannelURL == nil {
		url = "nil pointer"
	} else {
		url = *c.ChannelURL
	}

	return fmt.Sprintf("(tg_id: %d | channel_name: %s | ChannelURL: %s | Status: %s)",
		c.TelegramID, c.ChannelName, url, c.Status)
}
