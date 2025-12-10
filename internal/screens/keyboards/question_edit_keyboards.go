package keyboards

import (
	"strconv"
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetQuestionListKeyboard - клавиатура для списка вопросов с кнопками выбора
func GetQuestionListKeyboard(userState *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	checklistData, ok := userState.CurrentCheckList.(*types.SimpleCheckList)
	if !ok {
		return GetBackKeyboard(userState)
	}

	total := len(checklistData.Questions)
	if total == 0 {
		return GetBackKeyboard(userState)
	}

	// Создаем кнопки с карандашом и номером
	var rows [][]tgbotapi.KeyboardButton

	// Добавляем кнопки вопросов по 3 в ряд
	for i := 0; i < total; i++ {
		if i%3 == 0 {
			// Начинаем новую строку
			row := []tgbotapi.KeyboardButton{}

			// Добавляем до 3 кнопок в текущую строку
			for j := 0; j < 3 && i+j < total; j++ {
				questionNumber := i + j + 1
				// Создаем кнопку с карандашом и номером: ✏️ 1, ✏️ 2 и т.д.
				buttonText := "✏️ " + strconv.Itoa(questionNumber)
				row = append(row, tgbotapi.NewKeyboardButton(buttonText))
			}

			rows = append(rows, row)
		}
	}

	// Кнопка возврата
	rows = append(rows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(buttons.BtnBack),
	})

	return tgbotapi.NewReplyKeyboard(rows...)
}

// GetQuestionDetailKeyboard - клавиатура для деталей вопроса
func GetQuestionDetailKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnEditQuestionText),
			tgbotapi.NewKeyboardButton(buttons.BtnEditQuestionType),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnDeleteQuestion),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnBack),
		),
	)
}
