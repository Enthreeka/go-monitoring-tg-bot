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
)

type Channel struct {
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

	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		store:                  make(map[int64]*StoreData, 15),
		totalSuccessfulSentMsg: make(map[int64]channelStat),
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

func (s *Store) IncrementSuccessfulSentMsg(channelID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	today := time.Now().Day()

	if d, ok := s.totalSuccessfulSentMsg[channelID]; !ok {

		s.totalSuccessfulSentMsg[channelID] = channelStat{
			day:       today,
			countSend: 1,
		}

	} else {

		if d.day != today {
			s.totalSuccessfulSentMsg[channelID] = channelStat{
				day:       today,
				countSend: 1,
			}

		} else {
			s.totalSuccessfulSentMsg[channelID] = channelStat{
				day:       today,
				countSend: d.countSend + 1,
			}
		}

	}
}

func (s *Store) GetSuccessfulSentMsg(channelID int64) (int, int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	day := time.Now().Day()
	d, ok := s.totalSuccessfulSentMsg[channelID]
	if !ok {
		return day, 0
	}
	return d.day, d.countSend
}
