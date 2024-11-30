package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.GetAccessToken(chatID)
}

func (b *Bot) initAuth(chatID int64) (string, error) {
	ctx := context.Background()

	requestToken, err := b.pocketClient.GetRequestToken(ctx, b.redirectUrl)
	if err != nil {
		return "", err
	}

	err = b.tokenRepository.SaveRequestToken(chatID, requestToken)
	if err != nil {
		return "", err
	}

	redirectUrl, err := b.pocketClient.GetRedirectUrl(requestToken, b.redirectUrl)
	if err != nil {
		return "", err
	}

	newMsg := tgbotapi.NewMessage(chatID, redirectUrl)
	redirectUrl, err = b.generateRedirectURL(newMsg.Text, chatID)
	if err != nil {
		return "", err
	}
	return redirectUrl, nil
}
