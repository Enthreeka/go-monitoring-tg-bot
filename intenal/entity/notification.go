package entity

import "fmt"

type Notification struct {
	ID               int     `json:"id"`
	ChannelID        int64   `json:"channel_id"`
	NotificationText string  `json:"notification_text"`
	FileID           *string `json:"file_id"`
	FileType         *string `json:"file_type"`
	ButtonURL        *string `json:"button_url"`
}

func (n Notification) String() string {
	return fmt.Sprintf("(id: %d | channel_id: %d | notification_text: %s | file_id: %s | file_type: %s | button_url: %s)",
		n.ID, n.ChannelID, n.NotificationText, *n.FileID, *n.FileType, *n.ButtonURL)
}
