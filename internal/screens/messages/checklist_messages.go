package messages

import (
	"fmt"
	"telegram-bot/internal/screens/formatters"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
)

// CreateSimpleChecklistNameMessage - сообщение для создания названия простого чек-листа
func CreateSimpleChecklistNameMessage(state *state.UserState) string {
	return "Введите название чек-листа:"
}

// CreateBlockChecklistNameMessage - сообщение для создания названия чек-листа с блоками
func CreateBlockChecklistNameMessage(state *state.UserState) string {
	return "Введите название чек-листа с блоками:"
}

// SimpleChecklistEditorMessage - сообщение для редактора простого чек-листа
func SimpleChecklistEditorMessage(state *state.UserState) string {
	if state != nil && state.HasCheckList() {
		return "Редактор чек-листа - " + state.CurrentCheckList.GetName()
	}
	return "Редактор чек-листа"
}

// BlockChecklistEditorMessage - сообщение для редактора чек-листа с блоками
func BlockChecklistEditorMessage(state *state.UserState) string {
	if state != nil && state.HasCheckList() {
		checklistData, ok := state.CurrentCheckList.(*types.BlockedCheckList)
		if ok {
			totalBlocks := len(checklistData.Blocks)
			totalQuestions := 0
			for _, block := range checklistData.Blocks {
				totalQuestions += len(block.Questions)
			}

			return fmt.Sprintf("🧱 Редактор чек-листа: %s\n\n📊 Статистика:\n• Блоков: %d\n• Вопросов: %d\n\nВыберите блок для редактирования или добавьте новый:",
				checklistData.Name, totalBlocks, totalQuestions)
		}
		return "Редактор чек-листа с блоками - " + state.CurrentCheckList.GetName()
	}
	return "Редактор чек-листа с блоками"
}

// EditChecklistTitleMessage - сообщение для редактирования названия чек-листа
func EditChecklistTitleMessage(state *state.UserState) string {
	currentName := ""
	if state != nil && state.CurrentCheckList != nil {
		currentName = state.CurrentCheckList.GetName()
	}

	return "✏️ Редактирование названия чек-листа\n\n" +
		"Текущее название: " + currentName + "\n\n" +
		"Введите новое название:"
}

// ChecklistPreviewMessage - сообщение для превью чек-листа
func ChecklistPreviewMessage(state *state.UserState) string {
	if state == nil || !state.HasCheckList() {
		return "❌ Чек-лист не найден"
	}

	return formatters.FormatChecklistPreview(state)
}

// ConfirmExitToMainMenuMessage - сообщение для подтверждения выхода
func ConfirmExitToMainMenuMessage(state *state.UserState) string {
	return "⚠️ ВНИМАНИЕ!\n\n" +
		"Вы собираетесь выйти из редактора чек-листа.\n\n" +
		"Все несохраненные изменения будут потеряны!\n\n" +
		"Вы уверены, что хотите выйти?"
}
