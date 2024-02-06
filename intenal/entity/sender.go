package entity

type Sender struct {
	ID                int    `json:"id"`
	ChannelTelegramID int64  `json:"channel_tg_id"`
	Message           string `json:"message"`

	ChannelName string `json:"channel_name,omitempty"`
}
