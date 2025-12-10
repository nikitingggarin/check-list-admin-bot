package messages

import (
	"fmt"
	"strings"
	"telegram-bot/internal/formatters"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/storage/models"
)

// MyChecklistsListMessage - сообщение для списка чек-листов
func MyChecklistsListMessage(state *state.UserState) string {
	checklists, ok := state.Data["my_checklists"].([]models.Checklist)
	if !ok || len(checklists) == 0 {
		return "📭 У вас пока нет черновиков чек-листов."
	}

	var sb strings.Builder
	sb.WriteString("📝 Ваши черновики чек-листов:\n\n")

	for i, checklist := range checklists {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, checklist.Name))
		sb.WriteString(fmt.Sprintf("   🏷️ ID: %d\n", checklist.ID))
		sb.WriteString(fmt.Sprintf("   📅 Создан: %s\n", checklist.CreatedAt.Format("02.01.2006 15:04")))
		sb.WriteString("\n")
	}

	sb.WriteString("Выберите чек-лист для работы:")

	return sb.String()
}

// ChecklistDetailMessage - сообщение для деталей чек-листа (ТЕПЕРЬ С ПРЕВЬЮ)
func ChecklistDetailMessage(state *state.UserState) string {
	checklist, ok := state.Data["current_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		return "❌ Чек-лист не найден"
	}

	questions, _ := state.Data["checklist_questions"].([]models.Question)
	answerOptions, _ := state.Data["checklist_answer_options"].([]models.AnswerOption)
	hasBlocks, _ := state.Data["has_blocks"].(bool)
	blocks, _ := state.Data["checklist_blocks"].([]models.QuestionBlock)
	templates, _ := state.Data["checklist_templates"].([]models.ChecklistTemplate)

	return formatters.FormatChecklist(checklist, hasBlocks, blocks, templates, questions, answerOptions)
}

// ConfirmDeleteChecklistMessage - сообщение для подтверждения удаления
func ConfirmDeleteChecklistMessage(state *state.UserState) string {
	checklist, ok := state.Data["current_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		return "❌ Чек-лист не найден"
	}

	return "🗑️ УДАЛЕНИЕ ЧЕК-ЛИСТА\n\n" +
		"Название: " + checklist.Name + "\n" +
		"ID: " + fmt.Sprintf("%d", checklist.ID) + "\n\n" +
		"⚠️ Внимание! Это действие нельзя отменить.\n" +
		"Все вопросы и ответы будут удалены.\n\n" +
		"Вы уверены, что хотите удалить этот чек-лист?"
}
