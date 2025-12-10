package keyboards

import (
	"strconv"
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/storage/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetPublishedChecklistsKeyboard - клавиатура для списка опубликованных/отмененных чек-листов
func GetPublishedChecklistsKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	// Получаем список чек-листов из состояния
	checklists, ok := state.Data["published_checklists"].([]models.Checklist)
	if !ok || len(checklists) == 0 {
		// Если нет чек-листов - только кнопка назад
		return GetBackKeyboard(state)
	}

	var rows [][]tgbotapi.KeyboardButton

	// Создаем кнопки с карандашом и номером по 3 в ряд
	for i := 0; i < len(checklists); i++ {
		if i%3 == 0 {
			row := []tgbotapi.KeyboardButton{}

			// Добавляем до 3 кнопок в текущую строку
			for j := 0; j < 3 && i+j < len(checklists); j++ {
				checklistNumber := i + j + 1
				// Создаем кнопку с карандашом и номером: ✏️ 1, ✏️ 2 и т.д.
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

// GetPublishedChecklistDetailKeyboard - клавиатура для деталей опубликованного/отмененного чек-листа
func GetPublishedChecklistDetailKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	// Определяем тип чек-листа (published/unpublished)
	checklistType, ok := state.Data["current_checklist_type"].(string)
	if !ok {
		// fallback
		return GetBackKeyboard(state)
	}

	if checklistType == "published" {
		// Для опубликованных: снять с публикации + назад
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnUnPublish),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnBack),
			),
		)
	} else {
		// Для отмененных: вернуть в публикацию + назад
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnPublish),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnBack),
			),
		)
	}
}
