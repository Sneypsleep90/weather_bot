package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math"
	"time"
	"weatherr_bot/clients/openweather"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
	}

}

func (h *Handler) Start() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {
		h.HandleUpdate(update)
	}

}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return

	}

	if update.Message.IsCommand() && update.Message.Command() == "start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, sendHello)
		h.bot.Send(msg)
	}

	if update.Message.IsCommand() && update.Message.Command() == "time" {
		currentTime := time.Now().Format("15:04:05")
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Время: %s", currentTime))
		h.bot.Send(msg)

	}

	if !update.Message.IsCommand() {

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		coordinates, err := h.owClient.Coordinates(update.Message.Text)
		if err != nil {
			log.Println(err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить город!")
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}
		weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
		if err != nil {
			log.Println(err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%w", err))
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Погода в %s %d градусов", update.Message.Text, int(math.Round(weather.Temp))),
		)
		msg.ReplyToMessageID = update.Message.MessageID

		h.bot.Send(msg)

	}

}
