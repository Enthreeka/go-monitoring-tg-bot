package stateful

import (
	"sync"
	"time"
)

const (
	OperationUpdateText   = "update_text"
	OperationUpdateFile   = "update_file"
	OperationUpdateButton = "update_button"

	OperationAddAdmin      = "admin"
	OperationAddSuperAdmin = "superAdmin"
	OperationDeleteAdmin   = "user"

	OperationAddBot    = "add_bot"
	OperationDeleteBot = "delete_bot"

	OperationSetTimer = "set_timer"

	OperationUpdateQuestion = "update_question"
)

type Channel struct {
	MessageID   int
	ChannelName string

	OperationType string
}

type Sender struct {
	MessageID   int
	ChannelName string
}

type Notification struct {
	MessageID   int
	ChannelName string

	OperationType string
}

type Admin struct {
	MessageID int

	OperationType string
}

type SpamBot struct {
	MessageID int

	OperationType string
}

type StoreData struct {
	Admin        *Admin        `json:"admin"`
	Sender       *Sender       `json:"sender"`
	Channel      *Channel      `json:"channel"`
	Notification *Notification `json:"notification"`
	SpamBot      *SpamBot      `json:"spamBot"`
}

type Store struct {
	store                  map[int64]*StoreData
	totalSuccessfulSentMsg map[int64]channelStat

	userCaptcha map[int64]Captcha

	mu sync.RWMutex
}

type Captcha struct {
	ChannelName string
	Expire      time.Time
}

func NewStore() *Store {
	return &Store{
		store:                  make(map[int64]*StoreData, 15),
		totalSuccessfulSentMsg: make(map[int64]channelStat),
		userCaptcha:            make(map[int64]Captcha, 100),
	}
}

func (s *Store) Set(data *StoreData, userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[userID] = data
}

func (s *Store) Read(userID int64) (*StoreData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.store[userID]
	if !ok {
		return nil, false
	}

	return d, true
}

func (s *Store) Delete(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, userID)
}

type channelStat struct {
	countSend int64
	day       int
}
