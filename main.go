package main

import (
	"context"
	"log"
	"os"
	"weatherr_bot/clients/handler"
	"weatherr_bot/clients/openweather"
	"weatherr_bot/repo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := godotenv.Load()

	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("error connecting to db", err)
	}
	defer conn.Close()

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("error ping db")
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	usesRepo := repo.New(conn)

	owClient := openweather.New(os.Getenv("OPENWEATHERAPI_KEY"))
	botHandler := handler.New(bot, owClient, usesRepo)

	botHandler.Start()

}
