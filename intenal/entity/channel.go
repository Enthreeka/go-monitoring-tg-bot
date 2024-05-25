package entity

import (
	"fmt"
	"strconv"
	"strings"
)

type Channel struct {
	ID          int     `json:"id"`
	TelegramID  int64   `json:"tg_id"`
	ChannelName string  `json:"channel_name"`
	ChannelURL  *string `json:"channel_url"`
	AcceptTimer int     `json:"accept_timer"`

	Status       string `json:"status"`
	WaitingCount int    `json:"waiting_count,omitempty"`
	NeedCaptcha  bool   `json:"need_captcha"`
}

func (c Channel) String() string {
	var url string
	if c.ChannelURL == nil {
		url = "nil pointer"
	} else {
		url = *c.ChannelURL
	}

	return fmt.Sprintf("(tg_id: %d | channel_name: %s | ChannelURL: %s | Status: %s | need_captcha: %v)",
		c.TelegramID, c.ChannelName, url, c.Status, c.NeedCaptcha)
}

func GetID(data string) int {
	parts := strings.Split(data, "_")
	if len(parts) > 3 {
		return 0
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}

	return id
}
