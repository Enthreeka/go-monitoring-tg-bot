package entity

import "fmt"

type SpamBot struct {
	ID      int    `json:"id"`
	Token   string `json:"token"`
	BotName string `json:"bot_name"`
}

func (s *SpamBot) String() string {
	return fmt.Sprintf("(id: %d | token: %s | channel_name: %s)", s.ID, s.Token, s.BotName)
}
