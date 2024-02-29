package spam

import (
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type spammerBots struct {
	storageBot map[string]*tgbotapi.BotAPI

	log *logger.Logger
	mu  sync.RWMutex
}

type SpamBot interface {
	InitializeBot(token string) (string, error)
	Read(botName string) (*tgbotapi.BotAPI, bool)
	Delete(botName string)
	Range(f func(key, value any) error) error
}

func NewSpammerBot(log *logger.Logger) SpamBot {
	return &spammerBots{
		storageBot: make(map[string]*tgbotapi.BotAPI, 20),
		log:        log,
	}
}

func (s *spammerBots) InitializeBot(token string) (string, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		s.log.Error("failed to load token %v", err)
	}

	s.log.Info("Authorized on account %s", bot.Self.UserName)

	s.mu.Lock()
	s.storageBot[bot.Self.UserName] = bot
	s.mu.Unlock()

	return bot.Self.UserName, nil
}

func (s *spammerBots) Read(botName string) (*tgbotapi.BotAPI, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	d, ok := s.storageBot[botName]
	if !ok {
		return nil, false
	}

	return d, true
}

func (s *spammerBots) Delete(botName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.storageBot, botName)
}

func (s *spammerBots) Range(f func(key, value any) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for key, user := range s.storageBot {
		if err := f(key, user); err != nil {
			return err
		}
	}

	return nil
}
