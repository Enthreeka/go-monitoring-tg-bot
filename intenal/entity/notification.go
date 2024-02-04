package entity

import (
	"fmt"
)

type Notification struct {
	ID               int     `json:"id"`
	ChannelID        int64   `json:"channel_id"`
	NotificationText *string `json:"notification_text"`
	FileID           *string `json:"file_id"`
	FileType         *string `json:"file_type"`
	ButtonURL        *string `json:"button_url"`

	ChannelName string `json:"channel_name,omitempty"`
}

func (n Notification) String() string {
	var (
		text      string
		fileID    string
		fileType  string
		buttonURL string
	)

	if n.NotificationText == nil {
		text = "nil"
	} else {
		text = *n.NotificationText
	}

	if n.FileID == nil {
		fileID = "nil"
	} else {
		fileID = *n.FileID
	}

	if n.FileType == nil {
		fileType = "nil"
	} else {
		fileType = *n.FileType
	}

	if n.ButtonURL == nil {
		buttonURL = "nil"
	} else {
		buttonURL = *n.ButtonURL
	}

	return fmt.Sprintf("(id: %d | channel_id: %d | notification_text: %s | file_id: %s | file_type: %s | button_url: %s)",
		n.ID, n.ChannelID, text, fileID, fileType, buttonURL)
}
