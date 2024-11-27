package main

import (
	"github.com/StandyBee/gotgbot/pkg/telegram"
	pocket "github.com/StandyBee/pocketSDK"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	pcktClient, err := pocket.NewClient(os.Getenv("POCKET_CONSUMER_KEY"))
	if err != nil {
		log.Panic(err)
	}

	telegramBot := telegram.NewBot(bot, pcktClient, "https://google.com")
	err = telegramBot.Start()
	if err != nil {
		log.Fatal(err)
	}
}
