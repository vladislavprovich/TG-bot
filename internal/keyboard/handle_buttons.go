package keyboard

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/internal/service"
)

// todo redis ?
var userStates = make(map[int64]*models.UserAction)

type handleButtons struct {
	service service.UrlService
	logger  *logrus.Logger
}

func (h *handleButtons) HandleCallbackQuery(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	query := update.CallbackQuery
	userID := query.From.ID

	switch query.Data {
	case "create_short_url":
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "What you want?")
		msg.ReplyMarkup = CreateURL()
		bot.Send(msg)
	case "rand_url":
		// add action to tap buttons
		userStates[userID] = &models.UserAction{Action: "rand_url"}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Send only original url:")
		bot.Send(msg)

	case "cust_url":
		// add action to tap buttons
		userStates[userID] = &models.UserAction{Action: "cust_url"}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Send original url\n after short url:")
		bot.Send(msg)

	case "list_short_urls":
		urls, err := h.service.GetListUrl(ctx, models.GetListRequest{TgID: userID})
		if err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error retrieving URL list.")
			bot.Send(msg)
			return
		}

		urlList := "Your URLs:\n"
		if len(urls) == 0 {
			urlList += "No URLs found.\n"
		} else {
			for _, url := range urls {
				urlList += fmt.Sprintf("Original: %s\nShort: %s\n\n", url.OriginalUrl, url.ShortUrl)
			}
		}

		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, urlList)
		msg.ReplyMarkup = DeleteShortURL()
		bot.Send(msg)

	case "settings":
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Settings:")
		msg.ReplyMarkup = ClearAndBack()
		bot.Send(msg)

	case "back_to_main":
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Choose action in menu.")
		msg.ReplyMarkup = MainMenu()
		bot.Send(msg)

	case "clear_history":
		err := h.service.DeleteAllUrl(ctx, models.DeleteAllUrl{TgID: userID})
		if err != nil {
			msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Error history del.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "history del.")
		bot.Send(msg)
	case "delete_short_url":
		err := h.service.DeleteShortUrl(ctx, models.DeleteShortUrl{TgID: userID})
		if err != nil {
		}
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "Url deleted successfully.")
		bot.Send(msg)
	}

}

func (h *handleButtons) HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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
	case "rand_url":
		userStates[userID] = &models.UserAction{Action: "rand_url"}
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

	case "cust_url":
		userStates[userID] = &models.UserAction{Action: "cust_url"}
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
