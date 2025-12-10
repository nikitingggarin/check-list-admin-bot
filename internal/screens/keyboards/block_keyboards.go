package keyboards

import (
	"fmt"
	"strconv"
	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetChecklistBlockEditorKeyboard - основная клавиатура редактора блоков (показывает блоки + действия)
func GetChecklistBlockEditorKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	// Если нет чек-листа или он не блокированный - возвращаем базовую клавиатуру
	if state == nil || !state.HasCheckList() {
		return getBaseBlockActionsKeyboard()
	}

	checklistData, ok := state.CurrentCheckList.(*types.BlockedCheckList)
	if !ok {
		return getBaseBlockActionsKeyboard()
	}

	// Если нет блоков - показываем только кнопки действий
	if len(checklistData.Blocks) == 0 {
		return getBaseBlockActionsKeyboard()
	}

	var rows [][]tgbotapi.KeyboardButton

	// Добавляем кнопки блоков
	for i, block := range checklistData.Blocks {
		// Форматируем текст кнопки блока
		buttonText := formatBlockButton(&block, i+1)
		row := []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton(buttonText)}
		rows = append(rows, row)
	}

	// Добавляем разделительную пустую строку
	rows = append(rows, []tgbotapi.KeyboardButton{})

	// Добавляем кнопки действий
	rows = append(rows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(buttons.BtnAddBlock),
		tgbotapi.NewKeyboardButton(buttons.BtnPreview),
	})

	// Добавляем кнопку редактирования названия чек-листа
	rows = append(rows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(buttons.BtnEditTitleChecklist),
	})

	rows = append(rows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(buttons.BtnBackToMainMenu),
	})

	return tgbotapi.NewReplyKeyboard(rows...)
}

// GetBlockEditorKeyboard - клавиатура редактора конкретного блока
func GetBlockEditorKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	// Проверяем, есть ли вопросы в текущем блоке
	hasQuestions := false
	if state != nil && state.HasCheckList() {
		if checklistData, ok := state.CurrentCheckList.(*types.BlockedCheckList); ok {
			if blockIdx, ok := state.Data["current_block_index"].(int); ok {
				if blockIdx >= 0 && blockIdx < len(checklistData.Blocks) {
					hasQuestions = len(checklistData.Blocks[blockIdx].Questions) > 0
				}
			}
		}
	}

	if hasQuestions {
		// Есть вопросы - показываем все кнопки
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnAddQuestion),
				tgbotapi.NewKeyboardButton(buttons.BtnEditTitleBlockChecklist),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnEditQuestionChecklist),
				tgbotapi.NewKeyboardButton(buttons.BtnPreviewBlock),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnBackToBlockList),
			),
		)
	} else {
		// Нет вопросов - скрываем кнопки редактирования вопросов и превью блока
		return tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnAddQuestion),
				tgbotapi.NewKeyboardButton(buttons.BtnEditTitleBlockChecklist),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(buttons.BtnBackToBlockList),
			),
		)
	}
}

// GetEditBlockNameKeyboard - клавиатура для ввода/редактирования названия блока
func GetEditBlockNameKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnBack), //BtnCancel
		),
	)
}

// formatBlockButton - форматирует текст кнопки блока
func formatBlockButton(block *types.Block, number int) string {
	questionCount := len(block.Questions)
	var emoji string
	// Выбираем эмодзи в зависимости от количества вопросов
	switch {
	case questionCount == 0:
		emoji = "📭" // пустой блок
	case questionCount < 5:
		emoji = "🧱" // обычный блок
	default:
		emoji = "🏗️" // большой блок
	}
	name := block.Name
	return fmt.Sprintf("%s %d. %s (%d)", emoji, number, name, questionCount)
}

// getBaseBlockActionsKeyboard - базовая клавиатура действий (когда нет блоков)
func getBaseBlockActionsKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttons.BtnAddBlock),
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

// GetBlockQuestionsKeyboard - клавиатура для списка вопросов в блоке
func GetBlockQuestionsKeyboard(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	// Проверяем наличие чек-листа и блока
	if state == nil || !state.HasCheckList() {
		return GetBackKeyboard(state)
	}

	checklistData, ok := state.CurrentCheckList.(*types.BlockedCheckList)
	if !ok {
		return GetBackKeyboard(state)
	}

	// Получаем индекс текущего блока
	blockIdx, ok := state.Data["current_block_index"].(int)
	if !ok {
		return GetBackKeyboard(state)
	}

	if blockIdx < 0 || blockIdx >= len(checklistData.Blocks) {
		return GetBackKeyboard(state)
	}

	block := checklistData.Blocks[blockIdx]
	total := len(block.Questions)
	if total == 0 {
		return GetBackKeyboard(state)
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
