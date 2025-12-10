package messages

import (
	"fmt"
	"strings"
	"telegram-bot/internal/formatters"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/storage/models"
)

// PublishedChecklistsListMessage - сообщение для списка опубликованных/отмененных чек-листов
func PublishedChecklistsListMessage(state *state.UserState) string {
	checklists, ok := state.Data["published_checklists"].([]models.Checklist)
	if !ok || len(checklists) == 0 {
		return "📭 Список чек-листов пуст."
	}

	checklistType, _ := state.Data["checklists_type"].(string)

	var sb strings.Builder

	if checklistType == "published" {
		sb.WriteString("🚀 Ваши опубликованные чек-листы:\n\n")
	} else {
		sb.WriteString("🚫 Ваши отмененные публикации:\n\n")
	}

	for i, checklist := range checklists {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, checklist.Name))
		sb.WriteString(fmt.Sprintf("   🏷️ ID: %d\n", checklist.ID))

		createdAt := checklist.CreatedAt
		if !createdAt.IsZero() && createdAt.Year() > 1 {
			sb.WriteString(fmt.Sprintf("   📅 Создан: %s\n", createdAt.Format("02.01.2006 15:04")))
		} else {
			sb.WriteString("   📅 Создан: (дата недоступна)\n")
		}

		sb.WriteString("\n")
	}

	sb.WriteString("Выберите чек-лист для работы:")

	return sb.String()
}

// PublishedChecklistDetailMessage - сообщение для деталей опубликованного/отмененного чек-листа
func PublishedChecklistDetailMessage(state *state.UserState) string {
	checklist, ok := state.Data["current_published_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		return "❌ Чек-лист не найден"
	}

	questions, _ := state.Data["published_checklist_questions"].([]models.Question)
	answerOptions, _ := state.Data["published_checklist_answer_options"].([]models.AnswerOption)
	hasBlocks, _ := state.Data["published_has_blocks"].(bool)
	blocks, _ := state.Data["published_checklist_blocks"].([]models.QuestionBlock)
	templates, _ := state.Data["published_checklist_templates"].([]models.ChecklistTemplate)

	return formatters.FormatChecklist(checklist, hasBlocks, blocks, templates, questions, answerOptions)
}
