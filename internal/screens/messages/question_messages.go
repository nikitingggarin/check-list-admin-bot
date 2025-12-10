package messages

import (
	"fmt"
	"strings"
	"telegram-bot/internal/formatters"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
)

// SelectQuestionTypeMessage - сообщение для выбора типа вопроса
func SelectQuestionTypeMessage(state *state.UserState) string {
	return "Выберите тип вопроса:\n\n" +
		"✅ Соответствие - Да/Нет\n" +
		"🔘 Одиночный выбор - Один правильный вариант из нескольких\n" +
		"☑️ Множественный выбор - Несколько правильных вариантов\n" +
		"📝 Текстовый ответ - Свободный текст"
}

// EnterQuestionTextMessage - сообщение для ввода текста вопроса
func EnterQuestionTextMessage(state *state.UserState) string {
	return "Введите текст вопроса:"
}

// EnterAnswerOptionsMessage - сообщение для ввода вариантов ответов
func EnterAnswerOptionsMessage(state *state.UserState) string {
	return "Введите варианты ответов (каждый с новой строки):\n\n" +
		"Пример:\n" +
		"Вариант 1\n" +
		"Вариант 2\n" +
		"Вариант 3"
}

// SelectCorrectAnswersMessage - сообщение для выбора правильных ответов
func SelectCorrectAnswersMessage(state *state.UserState) string {
	count, _ := state.Data["answer_options_count"].(int)
	if count < 2 {
		count = 2
	}

	questionTypeStr, _ := state.Data["selected_question_type"].(string)
	questionType := types.QuestionCategory(questionTypeStr)

	baseMessage := "Введите номера правильных ответов"
	availableOptions := "Доступные варианты: 1-" + fmt.Sprintf("%d", count)

	switch questionType {
	case types.CategorySingleChoice:
		return baseMessage + " (например: 1):\n\n" +
			availableOptions + "\n" +
			"🔘 Одиночный выбор: нужен ровно 1 правильный ответ"

	case types.CategoryMultipleChoice:
		return baseMessage + " (например: 1,3):\n\n" +
			availableOptions + "\n" +
			"☑️ Множественный выбор: нужно минимум 2 правильных ответа"

	default:
		return baseMessage + " (например: 1 или 1,3):\n\n" +
			availableOptions
	}
}

// ViewQuestionMessage - сообщение для просмотра списка вопросов
func ViewQuestionMessage(state *state.UserState) string {
	checklistData, _ := state.CurrentCheckList.(*types.SimpleCheckList)
	total := len(checklistData.Questions)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 Список вопросов (%d шт.)\n\n", total))

	for i, question := range checklistData.Questions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, question.Text))
		sb.WriteString(fmt.Sprintf("   (Тип: %s)\n\n", formatters.FormatQuestionType(question.Category)))
	}

	sb.WriteString("Выберите вопрос для редактирования:")
	return sb.String()
}

// EditQuestionTextMessage - сообщение для редактирования текста вопроса
func EditQuestionTextMessage(state *state.UserState) string {
	checklistData, _ := state.CurrentCheckList.(*types.SimpleCheckList)
	idx, _ := state.Data["edit_question_index"].(int)
	currentText := checklistData.Questions[idx].Text

	return "✏️ Редактирование текста вопроса\n\n" +
		"Текущий текст: " + currentText + "\n\n" +
		"Введите новый текст:"
}

// EditQuestionTypeMessage - сообщение для изменения типа вопроса
func EditQuestionTypeMessage(state *state.UserState) string {
	checklistData, _ := state.CurrentCheckList.(*types.SimpleCheckList)
	idx, _ := state.Data["edit_question_index"].(int)
	currentType := checklistData.Questions[idx].Category

	return "✏️ Изменение типа вопроса\n\n" +
		"Текущий тип: " + formatters.FormatQuestionType(currentType) + "\n\n" +
		"Выберите новый тип:"
}

// ConfirmDeleteQuestionMessage - сообщение для подтверждения удаления вопроса
func ConfirmDeleteQuestionMessage(state *state.UserState) string {
	var questionText string

	isBlockQuestion, _ := state.Data["is_edit_block_questions"].(bool)

	if isBlockQuestion {
		if blockedChecklist, ok := state.CurrentCheckList.(*types.BlockedCheckList); ok && blockedChecklist != nil {
			blockIdx, _ := state.Data["current_block_index"].(int)
			questionIdx, _ := state.Data["edit_question_index"].(int)

			if blockIdx >= 0 && blockIdx < len(blockedChecklist.Blocks) {
				block := blockedChecklist.Blocks[blockIdx]
				if questionIdx >= 0 && questionIdx < len(block.Questions) {
					questionText = block.Questions[questionIdx].Text
				}
			}
		}
	} else {
		if simpleChecklist, ok := state.CurrentCheckList.(*types.SimpleCheckList); ok && simpleChecklist != nil {
			idx, _ := state.Data["edit_question_index"].(int)
			if idx >= 0 && idx < len(simpleChecklist.Questions) {
				questionText = simpleChecklist.Questions[idx].Text
			}
		}
	}

	if questionText == "" {
		questionText = "(текст вопроса недоступен)"
	}

	return "🗑️ Удаление вопроса\n\n" +
		"Текст вопроса: " + questionText + "\n\n" +
		"Вы уверены что хотите удалить этот вопрос?\n" +
		"Это действие нельзя отменить."
}

// EditQuestionDetailMessage - сообщение для деталей вопроса
func EditQuestionDetailMessage(state *state.UserState) string {
	isBlockQuestion, _ := state.Data["is_edit_block_questions"].(bool)

	if isBlockQuestion {
		checklistData, ok := state.CurrentCheckList.(*types.BlockedCheckList)
		if !ok {
			return "❌ Ошибка: неверный тип чек-листа для редактирования вопроса в блоке"
		}

		blockIdx, ok := state.Data["current_block_index"].(int)
		if !ok {
			return "❌ Ошибка: блок не выбран"
		}

		questionIdx, ok := state.Data["edit_question_index"].(int)
		if !ok {
			return "❌ Ошибка: вопрос не выбран"
		}

		if blockIdx < 0 || blockIdx >= len(checklistData.Blocks) {
			return "❌ Ошибка: неверный индекс блока"
		}

		block := checklistData.Blocks[blockIdx]

		if questionIdx < 0 || questionIdx >= len(block.Questions) {
			return "❌ Ошибка: неверный индекс вопроса"
		}

		question := block.Questions[questionIdx]

		message := fmt.Sprintf("🧱 Редактирование вопроса в блоке '%s'\n\n", block.Name)
		message += fmt.Sprintf("📝 Текст: %s\n", question.Text)
		message += fmt.Sprintf("🎯 Тип: %s\n", formatters.FormatQuestionType(question.Category))

		if len(question.AnswerOptions) > 0 {
			message += "\n📊 Варианты ответов:\n"
			for i, opt := range question.AnswerOptions {
				correctMark := " "
				if opt.IsCorrect {
					correctMark = "✅"
				}
				message += fmt.Sprintf("%s %d. %s\n", correctMark, i+1, opt.Text)
			}
		} else {
			message += "\n⚠️ Нет вариантов ответов"
		}

		message += "\n\nВыберите действие для редактирования:"
		return message
	} else {
		checklistData, ok := state.CurrentCheckList.(*types.SimpleCheckList)
		if !ok {
			return "❌ Ошибка: неверный тип чек-листа"
		}

		idx, ok := state.Data["edit_question_index"].(int)
		if !ok {
			return "❌ Ошибка: вопрос не выбран"
		}

		if idx < 0 || idx >= len(checklistData.Questions) {
			return "❌ Ошибка: неверный индекс вопроса"
		}

		question := checklistData.Questions[idx]

		message := fmt.Sprintf("✏️ Редактирование вопроса №%d/%d:\n\n", idx+1, len(checklistData.Questions))
		message += fmt.Sprintf("📝 Текст: %s\n", question.Text)
		message += fmt.Sprintf("🎯 Тип: %s\n", formatters.FormatQuestionType(question.Category))

		if len(question.AnswerOptions) > 0 {
			message += "\n📊 Варианты ответов:\n"
			for i, opt := range question.AnswerOptions {
				correctMark := " "
				if opt.IsCorrect {
					correctMark = "✅"
				}
				message += fmt.Sprintf("%s %d. %s\n", correctMark, i+1, opt.Text)
			}
		} else {
			message += "\n⚠️ Нет вариантов ответов"
		}

		message += "\n\nВыберите действие для редактирования:"
		return message
	}
}
