package telegram

import (
	"fmt"
	rediss "github.com/StandyBee/gotgbot/database/credis"
	pocket "github.com/StandyBee/pocketSDK"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository *rediss.TokenRepository
	redirectUrl     string
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, tr *rediss.TokenRepository, redirectUrl string) *Bot {
	return &Bot{
		bot:             bot,
		pocketClient:    client,
		tokenRepository: tr,
		redirectUrl:     redirectUrl,
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)
	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (b *Bot) generateRedirectURL(redirectUrl string, chatId int64) (string, error) {
	return fmt.Sprintf("%s?chat_id=%d", redirectUrl, chatId), nil
}
