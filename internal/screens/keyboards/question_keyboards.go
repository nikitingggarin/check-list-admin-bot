package keyboards

import (
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetQuestionTypeKeyboard - клавиатура выбора типа вопроса (все в один столбик)
func GetQuestionTypeKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnCompliance),
			tgbotapi.NewKeyboardButton(buttons.BtnSingleChoice),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnMultipleChoice),
			tgbotapi.NewKeyboardButton(buttons.BtnTextAnswer),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnBack),
		),
	)
}
