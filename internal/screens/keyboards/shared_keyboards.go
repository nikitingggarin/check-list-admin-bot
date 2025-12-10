package keyboards

import (
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetConfirmationKeyboard - общая клавиатура подтверждения
func GetConfirmationKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnYes),
			tgbotapi.NewKeyboardButton(buttons.BtnNo),
		),
	)
}

func GetBackKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnBack),
		),
	)
}
