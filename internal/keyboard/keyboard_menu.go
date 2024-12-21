package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vladislavprovich/TG-bot/internal/models"
	"github.com/vladislavprovich/TG-bot/pkg"
	"strconv"
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

		redactionOriginalUrl := pkg.OriginalInfo(url.OriginalUrl) //only original url. Example: "http://google.com/qwerty" --> "google.com"
		redactionShortUrl := pkg.ShortInfo(url.ShortUrl)          //only short url. Example: "http://host:1111/qwerty" --> "qwerty"

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(fmt.Sprintf("%s", redactionOriginalUrl), url.OriginalUrl),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", redactionShortUrl), "ignore"), // button no usage. TG Api block localhost
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

func CreateURLStatusButton(stats []*models.GetUrlStatusResponse) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, url := range stats {
		redactionShortUrl := pkg.ShortInfo(url.ShortUrl) //only short url. Example: "http://host:1111/qwerty" --> "qwerty"
		regi := strconv.Itoa(url.RedirectCount)

		datePart := url.CreatedAt.Format("2006-01-02")
		timePart := url.CreatedAt.Format("15:04")
		formatDate := datePart + "\n" + timePart

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", redactionShortUrl), "ignore"), // button no usage. TG Api block localhost
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", formatDate), "ignore"),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", regi), "ignore"),
		)
		rows = append(rows, row)
	}

	backRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Back", "back_to_main"),
	)
	rows = append(rows, backRow)

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}
