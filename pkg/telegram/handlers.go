package telegram

import (
	"context"
	"fmt"
	pocket "github.com/StandyBee/pocketSDK"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/url"
	"os"
	"strings"
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

		if update.Message.Sticker != nil {
			b.handleSticker(update.Message)

			continue
		}
		b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	log.Printf("[%s] %s", msg.From.UserName, msg.Text)

	if !isValidURL(msg.Text) {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, "enter valid link")
		b.bot.Send(newMsg)
		return
	}

	accessToken, err := b.tokenRepository.GetAccessToken(msg.Chat.ID)
	if err != nil {
		fmt.Printf("ACCESS OTVALILSYA %s", err.Error())
		return
	}

	if accessToken == "" {
		fmt.Printf("ACCESS PUSTOVAT %s", err.Error())
		return
	}

	addItemReq := pocket.AddItemRequest{
		Url:         msg.Text,
		Title:       "",
		Tags:        nil,
		ConsumerKey: os.Getenv("POCKET_CONSUMER_KEY"),
		AccessToken: accessToken,
	}

	err = b.pocketClient.Add(context.Background(), &addItemReq)
	if err != nil {
		fmt.Printf("fail to add link %s", err.Error())
		return
	}

	newMsg := tgbotapi.NewMessage(msg.Chat.ID, "DELO SDELANO")
	b.bot.Send(newMsg)
}

func (b *Bot) handleCommand(command *tgbotapi.Message) error {
	switch command.Command() {
	case "start":
		err := b.handleStartCommand(command)
		return err
	default:
		newMsg := tgbotapi.NewMessage(command.Chat.ID, "/start only")
		_, err := b.bot.Send(newMsg)
		return err
	}
}

func (b *Bot) handleStartCommand(command *tgbotapi.Message) error {
	_, err := b.getAccessToken(command.Chat.ID)
	if err == nil {
		stickerFileID := "CAACAgIAAxkBAAOeZ0mZp2QpCw-T0t-W9LMYW7Xrwd4AAoUQAAKzgklKCFC-NhD0vR02BA"

		sticker := tgbotapi.NewStickerShare(command.Chat.ID, stickerFileID)

		_, sendErr := b.bot.Send(sticker)
		if sendErr != nil {
			return sendErr
		}
		return nil
	}

	newMsg, err := b.initAuth(command.Chat.ID)
	if err != nil {
		return err
	}

	response := tgbotapi.NewMessage(command.Chat.ID, fmt.Sprintf("tap link %s", newMsg))

	_, err = b.bot.Send(response)
	return err
}

func (b *Bot) handleSticker(msg *tgbotapi.Message) {
	if msg.Sticker != nil {
		log.Printf("Received sticker with File ID: %s", msg.Sticker.FileID)
	}
}

func isValidURL(input string) bool {
	parsedURL, err := url.Parse(input)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	if !strings.HasPrefix(parsedURL.Scheme, "http") {
		return false
	}

	return true
}
