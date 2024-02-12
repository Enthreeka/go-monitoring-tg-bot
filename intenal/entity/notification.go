package entity

import (
	"fmt"
	"strings"
)

type Notification struct {
	ID               int     `json:"id"`
	ChannelID        int64   `json:"channel_id"`
	NotificationText *string `json:"notification_text"`
	FileID           *string `json:"file_id"`
	FileType         *string `json:"file_type"`
	ButtonURL        *string `json:"button_url"`
	ButtonText       *string `json:"button_text"`

	ChannelName string `json:"channel_name,omitempty"`
}

func (n Notification) String() string {
	var (
		text       string
		fileID     string
		fileType   string
		buttonURL  string
		buttonText string
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

	if n.ButtonText == nil {
		buttonText = "nil"
	} else {
		buttonText = *n.ButtonText
	}

	return fmt.Sprintf("(notification_text: %s | file_id: %s | file_type: %s | button_url: %s"+
		"| button_text: %s)", text, fileID, fileType, buttonURL, buttonText)
}

func GetButtonData(text string) (string, string) {
	parts := strings.Split(text, "|")
	if len(parts) != 2 {
		return "", ""
	}
	return strings.Trim(parts[0], " "), strings.Trim(parts[1], " ")
}
