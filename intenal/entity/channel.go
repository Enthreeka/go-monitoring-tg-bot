package entity

import "fmt"

type Channel struct {
	UserID        int64  `json:"user_id"`
	ID            int    `json:"id"`
	StatusRequest string `json:"status_request"`
	User          User   `json:"user"`
}

func (c Channel) String() string {
	return fmt.Sprintf("(user_id: %d | id: %d | status_request: %s | user: %v)", c.UserID, c.ID, c.StatusRequest, c.User)
}
