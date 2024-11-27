package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)
			if err != nil {
				return
			}
			continue
		}
		b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	log.Printf("[%s] %s", msg.From.UserName, msg.Text)

	newMsg := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)

	b.bot.Send(newMsg)
}

func (b *Bot) handleCommand(command *tgbotapi.Message) error {
	switch command.Command() {
	case "start":
		err := b.handleStartCommand(command)
		return err
	default:
		newMsg := tgbotapi.NewMessage(command.Chat.ID, "HUI TEBE")
		_, err := b.bot.Send(newMsg)
		return err
	}
}

func (b *Bot) handleStartCommand(command *tgbotapi.Message) error {
	newMsg, err := b.initAuth(command.Chat.ID)
	if err != nil {
		return err
	}

	response := tgbotapi.NewMessage(command.Chat.ID, newMsg)

	_, err = b.bot.Send(response)
	return err
}
