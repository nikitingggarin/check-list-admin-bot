package keyboards

import (
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetChecklistSimpleEditorKeyboard - клавиатура редактора простого чек-листа
func GetChecklistSimpleEditorKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	// Проверяем, есть ли вопросы в текущем чек-листе
	hasQuestions := false
	if state != nil && state.CurrentCheckList != nil {
		if simpleChecklist, ok := state.CurrentCheckList.(*types.SimpleCheckList); ok {
			hasQuestions = len(simpleChecklist.Questions) > 0
		}
	}

	if hasQuestions {
		// Есть вопросы - показываем все кнопки
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnAddQuestion),
				tgbotapi.NewKeyboardButton(buttons.BtnPreview),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnEditTitleChecklist),
				tgbotapi.NewKeyboardButton(buttons.BtnEditQuestionChecklist),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnBackToMainMenu),
			),
		)
	} else {
		// Нет вопросов - скрываем кнопку редактирования вопросов
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnAddQuestion),
				tgbotapi.NewKeyboardButton(buttons.BtnPreview),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnEditTitleChecklist),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnBackToMainMenu),
			),
		)
	}
}

// GetChecklistPreviewKeyboard - клавиатура превью чек-листа
func GetChecklistPreviewKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnSaveDraft),
			tgbotapi.NewKeyboardButton(buttons.BtnSavePublish),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnBack),
		),
	)
}
