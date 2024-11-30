package main

import (
	"github.com/StandyBee/gotgbot/database/credis"
	"github.com/StandyBee/gotgbot/pkg/server"
	"github.com/StandyBee/gotgbot/pkg/telegram"
	pocket "github.com/StandyBee/pocketSDK"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := os.Getenv("REDIS_DB")
	intDB, err := strconv.Atoi(db)
	if err != nil {
		log.Fatal("Error converting credis DB to int")
	}

	redisClient := credis.NewRedisClient(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"), intDB)
	repository := credis.NewTokenRepository(redisClient)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_API_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	pcktClient, err := pocket.NewClient(os.Getenv("POCKET_CONSUMER_KEY"))
	if err != nil {
		log.Panic(err)
	}

	telegramBot := telegram.NewBot(bot, pcktClient, repository, "http://localhost")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	authServer := server.NewAuthorizationServer(pcktClient, "https://t.me/pocket_fsk_bot", repository)
	if err := authServer.Start(); err != nil {
		log.Fatal(err)
	}
}
