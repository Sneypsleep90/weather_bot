package handler

import (
	"context"
	"fmt"
	"log"
	"math"
	"weatherr_bot/clients/openweather"
	"weatherr_bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserRepository interface {
	GetUserCity(ctx context.Context, userId int64) (string, error)
	CreateUser(ctx context.Context, userId int64) error
	UpdateCity(ctx context.Context, userId int64, city string) error
	GetUser(ctx context.Context, userId int64) (*models.User, error)
}
type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
	userRepo UserRepository
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient, userRepo UserRepository) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
		userRepo: userRepo,
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
		return
	}

	ctx := context.Background()

	if update.Message.IsCommand() {
		err := h.ensureUser(ctx, update)
		if err != nil {
			log.Println("error ensureUser: ", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "произшала ошибка")
			msg.ReplyToMessageID = update.Message.MessageID
			h.bot.Send(msg)
			return
		}

		switch update.Message.Command() {
		case "city":
			h.HandleSetCity(ctx, update)
			return
		case "weather":
			h.HandleSendWeather(ctx, update)
			return
		default:
			h.HandleUnknownCommand(update)
			return

		}

	}
}

func (h *Handler) ensureUser(ctx context.Context, update tgbotapi.Update) error {
	user, err := h.userRepo.GetUser(ctx, update.Message.From.ID)
	if err != nil {
		return fmt.Errorf("error userRepo.GetUser: %w", err)
	}
	if user == nil {
		err := h.userRepo.CreateUser(ctx, update.Message.From.ID)
		if err != nil {
			return fmt.Errorf("error userRepo.CreateUser: %w", err)
		}
	}
	return nil
}

func (h *Handler) HandleSetCity(ctx context.Context, update tgbotapi.Update) {
	city := update.Message.CommandArguments()
	if err := h.userRepo.CreateUser(ctx, update.Message.From.ID); err != nil {
		log.Println("error userRepo.CreateUser: ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "произшала ошибка")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Город %s сохранен", city))
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)

}

func (h *Handler) HandleSendWeather(ctx context.Context, update tgbotapi.Update) {
	city, err := h.userRepo.GetUserCity(ctx, update.Message.From.ID)
	if err != nil {
		log.Println("error userRepo.GetUserCity : ", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "произшала ошибка")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	if city == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "сначала установите город!")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	coordinates, err := h.owClient.Coordinates(city)
	if err != nil {
		log.Printf("error owClient.Coordinates : %w", err)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не смогли получить город!")
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}
	weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
	if err != nil {
		log.Printf("error owClient.Weather: %w", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%w", err))
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("Погода в %s %d градусов", city, int(math.Round(weather.Temp))),
	)
	msg.ReplyToMessageID = update.Message.MessageID

	h.bot.Send(msg)
}

func (h *Handler) HandleUnknownCommand(update tgbotapi.Update) {
	log.Printf("unknow command [%s]%s", update.Message.From.UserName, update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "пока такой команды нет")
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}
