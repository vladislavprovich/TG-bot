package keyboard

import (
	"context"
	"database/sql"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/repository"
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
)

// todo redis ?
var userStates = make(map[int64]*models.UserAction)

type HandleButtons struct {
	service service.UrlService
	logger  *logrus.Logger
	repo    *repository.UserRepository
}

func NewHandleButtons(service service.UrlService, logger *logrus.Logger) *HandleButtons {
	return &HandleButtons{service: service, logger: logger}
}

func NewMessageHandler(service service.UrlService, logger *logrus.Logger) *HandleButtons {
	return &HandleButtons{service: service, logger: logger}
}

func (h *HandleButtons) HandleCallbackQuery(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		h.logger.Info("CallbackQuery is nil")
		return
	}

	query := update.CallbackQuery
	tgID := query.From.ID

	checkUserReq := models.CreateNewUserRequest{
		TgID: tgID,
	}

	userID, err := h.service.CreateUserByTgID(ctx, checkUserReq)
	if err != nil {
		h.logger.Error(err)
	}

	if query.Message == nil {
		h.logger.Info("Message in CallbackQuery is nil")
		return
	}

	if strings.HasPrefix(query.Data, "delete_short_url:") {
		shortUrl := strings.TrimPrefix(query.Data, "delete_short_url:")
		h.logger.Infof("Received callback to delete short URL: %s", shortUrl)

		if err = h.service.DeleteShortUrl(ctx, models.DeleteShortUrl{TgID: tgID, UserID: userID, ShortUrl: shortUrl}); err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Failed to delete URL.")
			msg.ReplyMarkup = BackMenu()
			bot.Send(msg)
			return
		}

		urls, err := h.service.GetListUrl(ctx, models.GetListRequest{UserID: userID, TgID: tgID})
		if err != nil || len(urls) == 0 {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "All URLs have been deleted.")
			msg.ReplyMarkup = BackMenu()
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Your URLs:")
		msg.ReplyMarkup = CreateURLListWithDeleteButtons(urls)
		bot.Send(msg)
		return
	}

	switch query.Data {
	case CreateShortURL:
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Choose an action from the menu.")
		msg.ReplyMarkup = CreateURL()
		bot.Send(msg)

	case RandomURL:
		userStates[tgID] = &models.UserAction{Action: RandomURL}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Send original URL:")
		bot.Send(msg)

	case CustomURL:
		userStates[tgID] = &models.UserAction{Action: CustomURL}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Send original URL:")
		bot.Send(msg)

	case ListURL:
		urls, err := h.service.GetListUrl(ctx, models.GetListRequest{UserID: userID, TgID: tgID})
		if err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error retrieving URL list.")
			msg.ReplyMarkup = BackMenu()
			bot.Send(msg)
			return
		}

		if len(urls) == 0 {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "You have no short URLs.")
			msg.ReplyMarkup = BackMenu()
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Your URLs: ")
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
		if err := h.service.DeleteAllUrl(ctx, models.DeleteAllUrl{UserID: userID, TgID: tgID}); err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error clearing history.")
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "History cleared.")
		msg.ReplyMarkup = BackMenu()
		bot.Send(msg)
	//todo
	case "show_url_stats":
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "coming soon:")
		msg.ReplyMarkup = BackMenu()
		bot.Send(msg)
	}

}

func (h *HandleButtons) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		h.logger.Info("Message is nil, skipping update.")
		return
	}

	tgID := update.Message.From.ID
	messageText := update.Message.Text
	var userID string

	user, exists := userStates[tgID]
	if !exists {
		userStates[tgID] = &models.UserAction{}
		h.logger.Infof("Created new user state for userID: %d\n", tgID)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose an action from the menu.")
		msg.ReplyMarkup = MainMenu()
		if _, err := bot.Send(msg); err != nil {
			h.logger.Errorf("Failed to send default menu: %v", err)
		}
		return
	}

	switch user.Action {
	case Start:
		if messageText == Start {
			createUserReq := models.CreateNewUserRequest{
				TgID:     tgID,
				UserName: update.Message.From.UserName,
			}
			var err error
			userID, err = h.service.CreateUserByTgID(ctx, createUserReq)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					h.logger.Errorf("Failed to create user: %v", err)
					h.logger.Errorf("USER ID: %v", userID)
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the bot! Choose an action from the menu.")
			msg.ReplyMarkup = MainMenu()
			if _, err = bot.Send(msg); err != nil {
				h.logger.Errorf("Failed to send start message: %v", err)
			}
			return
		}
	case RandomURL:
		req := models.CreateShortUrlRequest{
			TgID:        tgID,
			UserID:      userID,
			OriginalUrl: messageText,
		}

		resp, err := h.service.CreateShortUrl(ctx, req)
		if err != nil {
			h.logger.Errorf("Failed to create short URL: %v", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
			if _, err = bot.Send(msg); err != nil {
				h.logger.Errorf("Failed to send error message: %v", err)
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your short URL: "+resp.ShortUrl)
		msg.ReplyMarkup = BackMenu()
		if _, err = bot.Send(msg); err != nil {
			h.logger.Errorf("Failed to send short URL: %v", err)
		}

		userStates[tgID] = &models.UserAction{}

	case CustomURL:
		if user.OriginalURL == "" {
			user.OriginalURL = messageText

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Send custom short URL:")
			if _, err := bot.Send(msg); err != nil {
				h.logger.Errorf("Failed to request custom short URL: %v", err)
			}
		} else {
			user.CustomUrl = messageText
			req := models.CreateShortUrlRequest{
				TgID:        tgID,
				UserID:      userID,
				OriginalUrl: user.OriginalURL,
				CustomAlias: &user.CustomUrl,
			}

			resp, err := h.service.CreateShortUrl(ctx, req)
			if err != nil {
				h.logger.Errorf("Failed to create custom short URL: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error: "+err.Error())
				if _, err = bot.Send(msg); err != nil {
					h.logger.Errorf("Failed to send error message: %v", err)
				}
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your short URL: "+resp.ShortUrl)
			msg.ReplyMarkup = BackMenu()
			if _, err = bot.Send(msg); err != nil {
				h.logger.Errorf("Failed to send short URL: %v", err)
			}

			userStates[tgID] = &models.UserAction{}
		}

	default:
		h.logger.Errorf("Unexpected user action: %v\n", user.Action)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose an action from the menu.")
		msg.ReplyMarkup = MainMenu()
		if _, err := bot.Send(msg); err != nil {
			h.logger.Errorf("Failed to send default menu: %v", err)
		}
	}
}
