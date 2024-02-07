package stateful

import "sync"

const (
	OperationUpdateText   = "update_text"
	OperationUpdateFile   = "update_file"
	OperationUpdateButton = "update_button"
)

type Channel struct {
}

type Sender struct {
	ChannelName string
	MessageID   int
}

type Notification struct {
	ChannelName string
	MessageID   int

	OperationType string
}

type StoreData struct {
	Sender       *Sender       `json:"sender"`
	Channel      *Channel      `json:"channel"`
	Notification *Notification `json:"notification"`
}

type Store struct {
	store map[int64]*StoreData

	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		store: make(map[int64]*StoreData, 15),
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
