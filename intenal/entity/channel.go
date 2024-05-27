package entity

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Channel struct {
	ID              int             `json:"id"`
	TelegramID      int64           `json:"tg_id"`
	ChannelName     string          `json:"channel_name"`
	ChannelURL      *string         `json:"channel_url"`
	AcceptTimer     int             `json:"accept_timer"`
	Question        json.RawMessage `json:"question"`
	QuestionEnabled bool            `json:"question_enabled"`

	Status       string `json:"status"`
	WaitingCount int    `json:"waiting_count,omitempty"`
	NeedCaptcha  bool   `json:"need_captcha"`
}

// QuestionModel - включает все поля для опроса юзеров
type QuestionModel struct {
	ChanelNameBase64 string   `json:"-"`
	Question         string   `json:"question"`
	Answer           []Answer `json:"answer"`
}

type Answer struct {
	ID              int    `json:"-"`
	AnswerVariation string `json:"answer_variation"`
	Url             string `json:"url"`
	TextResult      string `json:"text_result"`
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

func ExtractValues(input string) (string, string) {
	parts := strings.Split(input, "_")
	if len(parts) != 3 {
		return "", ""
	}
	return parts[1], parts[2]
}

func IsValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}
