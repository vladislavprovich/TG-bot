package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vladislavprovich/TG-bot/internal/models"
)

func MainMenu() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Create short URL", "create_short_url"),
			tgbotapi.NewInlineKeyboardButtonData("List of all short URLs", "list_short_urls"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Show short URL stats", "show_url_stats"),
			tgbotapi.NewInlineKeyboardButtonData("Settings", "settings"),
		),
	)
	return &keyboard
}

func BackMenu() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "back_to_main"),
		),
	)
	return &keyboard
}

func ClearAndBack() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Clear History", "clear_history"),
			tgbotapi.NewInlineKeyboardButtonData("Back", "back_to_main"),
		),
	)
	return &keyboard
}

func CreateURL() *tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Create Random", "rand_url"),

			tgbotapi.NewInlineKeyboardButtonData("Create Custom", "cust_url"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "back_to_main"),
		),
	)
	return &keyboard
}

func CreateURLListWithDeleteButtons(urls []*models.GetListResponse) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, url := range urls {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", url.OriginalUrl), "url"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", url.ShortUrl), "url"),
			tgbotapi.NewInlineKeyboardButtonData("Delete", fmt.Sprintf("delete_short_url:%s", url.ShortUrl)),
		)
		rows = append(rows, row)
	}

	backRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Back", "back_to_main"),
	)
	rows = append(rows, backRow)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}
