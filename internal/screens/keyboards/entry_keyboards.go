package keyboards

import (
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetAuthKeyboard возвращает клавиатуру авторизации
func GetAuthKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnAuth),
		),
	)
}

// GetAdminMenu возвращает главное меню администратора
func GetAdminMenu(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnCreateSimpleChecklist),
			tgbotapi.NewKeyboardButton(buttons.BtnCreateBlockChecklist),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnMyChecklists),
			tgbotapi.NewKeyboardButton(buttons.BtnPublished),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnCanceled),
			tgbotapi.NewKeyboardButton(buttons.BtnStatistics),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnLogout),
		),
	)
}
