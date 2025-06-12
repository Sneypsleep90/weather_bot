package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"weatherr_bot/clients/handler"
	"weatherr_bot/clients/openweather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	
	owClient := openweather.New(os.Getenv("OPENWEATHERAPI_KEY"))
	botHandler := handler.New(bot, owClient)

	botHandler.Start()

}
