package messages

import (
	"fmt"
	"strings"
	"telegram-bot/internal/formatters"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
)

// EditBlockNameMessage - сообщение для ввода/редактирования названия блока
func EditBlockNameMessage(state *state.UserState) string {
	if blockIdx, ok := state.Data["current_block_index"].(int); ok {
		checklistData, _ := state.CurrentCheckList.(*types.BlockedCheckList)
		if blockIdx >= 0 && blockIdx < len(checklistData.Blocks) {
			currentName := checklistData.Blocks[blockIdx].Name
			return fmt.Sprintf("✏️ Редактирование названия блока\n\nТекущее название: %s\n\nВведите новое название:", currentName)
		}
	}
	return "Введите название нового блока:"
}

// BlockEditorMessage - сообщение для редактора блока
func BlockEditorMessage(state *state.UserState) string {
	blockIdx := state.Data["current_block_index"].(int)
	checklistData := state.CurrentCheckList.(*types.BlockedCheckList)
	block := checklistData.Blocks[blockIdx]
	questionCount := len(block.Questions)

	var statusEmoji string
	if questionCount == 0 {
		statusEmoji = "📭"
	} else if questionCount < 5 {
		statusEmoji = "🧱"
	} else {
		statusEmoji = "🏗️"
	}

	return fmt.Sprintf("%s Редактор блока: %s\n\n📊 Вопросов: %d\n\nВыберите действие:",
		statusEmoji, block.Name, questionCount)
}

// ConfirmExitBlockChecklistMessage - сообщение для подтверждения выхода из редактора блоков
func ConfirmExitBlockChecklistMessage(state *state.UserState) string {
	return "⚠️ ВНИМАНИЕ!\n\n" +
		"Вы собираетесь выйти из редактора чек-листа с блоками.\n\n" +
		"Все несохраненные изменения будут потеряны!\n\n" +
		"Вы уверены, что хотите выйти?"
}

// BlockViewQuestionsMessage - сообщение для списка вопросов в блоке
func BlockViewQuestionsMessage(state *state.UserState) string {
	checklistData := state.CurrentCheckList.(*types.BlockedCheckList)
	blockIdx := state.Data["current_block_index"].(int)
	block := checklistData.Blocks[blockIdx]
	total := len(block.Questions)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 Список вопросов в блоке '%s' (%d шт.)\n\n", block.Name, total))

	for i, question := range block.Questions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, question.Text))
		sb.WriteString(fmt.Sprintf("   (Тип: %s)\n\n", formatters.FormatQuestionType(question.Category)))
	}

	sb.WriteString("Выберите вопрос для редактирования:")
	return sb.String()
}
