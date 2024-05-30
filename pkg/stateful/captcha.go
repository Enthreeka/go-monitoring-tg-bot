package stateful

import (
	"context"
	"log"
	"time"
)

func (s *Store) Worker(ctx context.Context, minute int) {
	tick := time.NewTicker(time.Duration(minute) * time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped")
			return
		case <-tick.C:
			err := s.RangeCaptcha(func(key int64, value Captcha) error {
				if time.Now().After(value.Expire) {
					s.DeleteCaptcha(key)
				}
				return nil
			})
			if err != nil {
				log.Printf("Worker: %v", err)
			}
		}
	}
}

func (s *Store) SetCaptcha(data Captcha, userID int64) {
	s.mu.Lock()
	s.userCaptcha[userID] = data
	s.mu.Unlock()
}

func (s *Store) ReadCaptcha(userID int64) (Captcha, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.userCaptcha[userID]
	if !ok {
		return Captcha{}, false
	}

	return d, true
}

func (s *Store) DeleteCaptcha(userID int64) {
	s.mu.Lock()
	delete(s.userCaptcha, userID)
	s.mu.Unlock()
}

func (s *Store) RangeCaptcha(f func(key int64, value Captcha) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for key, user := range s.userCaptcha {
		if err := f(key, user); err != nil {
			return err
		}
	}

	return nil
}
