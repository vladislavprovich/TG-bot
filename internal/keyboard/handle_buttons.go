package keyboard

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/service"
	"strings"
)

const (
	Start          = "/start"
	CreateShortURL = "create_short_url"
	RandomURL      = "rand_url"
	CustomURL      = "cust_url"
	ListURL        = "list_short_urls"
	Settings       = "settings"
	BackMainMenu   = "back_to_main"
	ClearHistory   = "clear_history"
	DeleteShort    = "delete_short_url"
)

// todo redis ?
var userStates = make(map[int64]*models.UserAction)

type HandleButtons struct {
	service service.UrlService
	logger  *logrus.Logger
}

func NewHandleButtons(service service.UrlService, logger *logrus.Logger) *HandleButtons {
	return &HandleButtons{service: service, logger: logger}
}

func NewMessageHandler(service service.UrlService, logger *logrus.Logger) *HandleButtons {
	return &HandleButtons{service: service, logger: logger}
}

func (h *HandleButtons) HandleCallbackQuery(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	query := update.CallbackQuery
	userID := query.From.ID
	tgID := update.Message.From.ID
	username := update.Message.From.UserName

	if update.Message != nil && update.Message.Text == Start {
		err := h.service.CreateUserByTgID(ctx, models.CreateNewUserRequest{TgID: tgID, UserName: username})
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Error to init new user."))
			return
		}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Hello! Take options:")
		msg.ReplyMarkup = MainMenu()
		bot.Send(msg)
	}

	switch query.Data {
	case CreateShortURL:
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "What you want?")
		msg.ReplyMarkup = CreateURL()
		bot.Send(msg)
	case RandomURL:
		// add action to tap buttons
		userStates[userID] = &models.UserAction{Action: RandomURL}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Send only original url:")
		bot.Send(msg)

	case CustomURL:
		// add action to tap buttons
		userStates[userID] = &models.UserAction{Action: CustomURL}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Send original url\n after short url:")
		bot.Send(msg)

	case ListURL:
		urls, err := h.service.GetListUrl(ctx, models.GetListRequest{TgID: userID})
		if err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error retrieving URL list.")
			bot.Send(msg)
			return
		}

		if len(urls) == 0 {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "You have no short URLs.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Your URLs:")
		msg.ReplyMarkup = CreateURLListWithDeleteButtons(urls)
		bot.Send(msg)

	case Settings:
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Settings:")
		msg.ReplyMarkup = ClearAndBack()
		bot.Send(msg)

	case BackMainMenu:
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Choose action in menu.")
		msg.ReplyMarkup = MainMenu()
		bot.Send(msg)

	case ClearHistory:
		err := h.service.DeleteAllUrl(ctx, models.DeleteAllUrl{TgID: userID})
		if err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error history del.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "history del.")
		bot.Send(msg)
	case DeleteShort:
		shortUrl := strings.TrimPrefix(query.Data, DeleteShort)

		err := h.service.DeleteShortUrl(ctx, models.DeleteShortUrl{TgID: userID, ShortUrl: shortUrl})
		if err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Failed to delete URL.")
			bot.Send(msg)
			return
		}

		// update list if we del 1 url
		urls, err := h.service.GetListUrl(ctx, models.GetListRequest{TgID: userID})
		if err != nil || len(urls) == 0 {
			// if list all del
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "All URLs have been deleted.")
			bot.Send(msg)
			return
		}

		// New list
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Your URLs:\n")
		msg.ReplyMarkup = CreateURLListWithDeleteButtons(urls)
		bot.Send(msg)
	}

}

func (h *HandleButtons) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	user, exists := userStates[userID]
	query := update.CallbackQuery

	if !exists {
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Оберіть дію з меню.")
		msg.ReplyMarkup = MainMenu()
		bot.Send(msg)
		return
	}

	switch user.Action {
	case RandomURL:
		userStates[userID] = &models.UserAction{Action: RandomURL}
		req := models.CreateShortUrlRequest{
			OriginalUrl: update.Message.Text,
		}
		resp, err := h.service.CreateShortUrl(ctx, req)
		if err != nil {
			bot.Send(tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error: "+err.Error()))
			return
		}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Your short url: "+resp.ShortUrl)
		msg.ReplyMarkup = BackMenu()
		bot.Send(msg)

	case CustomURL:
		userStates[userID] = &models.UserAction{Action: CustomURL}
		if user.OriginalURL == "" {
			user.OriginalURL = update.Message.Text
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Send custom short URL:")
			bot.Send(msg)
		} else {
			user.CustomUrl = update.Message.Text
			req := models.CreateShortUrlRequest{
				OriginalUrl: user.OriginalURL,
				CustomAlias: &user.CustomUrl,
			}
			resp, err := h.service.CreateShortUrl(ctx, req)
			if err != nil {
				bot.Send(tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error: "+err.Error()))
				return
			}
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Your short url: "+resp.ShortUrl)
			msg.ReplyMarkup = BackMenu()
			bot.Send(msg)
		}
	}
}
