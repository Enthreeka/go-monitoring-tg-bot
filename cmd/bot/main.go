package main

import (
	"github.com/Entreeka/monitoring-tg-bot/intenal/bot"
	"github.com/Entreeka/monitoring-tg-bot/intenal/config"
	"github.com/Entreeka/monitoring-tg-bot/pkg/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed load config: %v", err)
	}

	if err := bot.Run(log, cfg); err != nil {
		log.Fatal("failed to run tgbot: %v", err)
	}
}
