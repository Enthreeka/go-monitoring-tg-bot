package entity

import "fmt"

type Request struct {
	ID            int    `json:"id"`
	UserID        int64  `json:"user_id"`
	StatusRequest string `json:"status_request"`
}

func (r Request) String() string {
	return fmt.Sprintf("(id: %d | user_id: %d | status_request: %s)", r.ID, r.UserID, r.StatusRequest)
}
