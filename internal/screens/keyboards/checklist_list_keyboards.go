package keyboards

import (
	"strconv"
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/storage/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetMyChecklistsKeyboard - клавиатура для списка чек-листов с кнопками-номерами (как у вопросов)
func GetMyChecklistsKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	// Получаем список чек-листов из состояния
	checklists, ok := state.Data["my_checklists"].([]models.Checklist)
	if !ok || len(checklists) == 0 {
		// Если нет чек-листов - только кнопка назад
		return GetBackKeyboard(state)
	}

	var rows [][]tgbotapi.KeyboardButton

	// Создаем кнопки с карандашом и номером по 3 в ряд (как у вопросов)
	for i := 0; i < len(checklists); i++ {
		if i%3 == 0 {
			row := []tgbotapi.KeyboardButton{}

			// Добавляем до 3 кнопок в текущую строку
			for j := 0; j < 3 && i+j < len(checklists); j++ {
				checklistNumber := i + j + 1
				// Создаем кнопку с карандашом и номером: ✏️ 1, ✏️ 2 и т.д. (как у вопросов)
				buttonText := "✏️ " + strconv.Itoa(checklistNumber)
				row = append(row, tgbotapi.NewKeyboardButton(buttonText))
			}

			rows = append(rows, row)
		}
	}

	// Добавляем разделительную строку
	rows = append(rows, []tgbotapi.KeyboardButton{})

	// Кнопка возврата
	rows = append(rows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(buttons.BtnBack),
	})

	return tgbotapi.NewReplyKeyboard(rows...)
}

// GetChecklistDetailKeyboard - клавиатура для деталей чек-листа
func GetChecklistDetailKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnEditCheckList),
			tgbotapi.NewKeyboardButton(buttons.BtnPublishChecklist),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnDeleteCheckList),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnBack),
		),
	)
}
