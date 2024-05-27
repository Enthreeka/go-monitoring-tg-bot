package stateful

import "time"

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
