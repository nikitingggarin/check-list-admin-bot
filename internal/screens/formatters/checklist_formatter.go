package formatters

import (
	"telegram-bot/internal/formatters"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
)

// FormatChecklistPreview форматирует чек-лист для превью
func FormatChecklistPreview(state *state.UserState) string {
	if state == nil || !state.HasCheckList() {
		return "❌ Чек-лист не найден"
	}

	checklistData := state.CurrentCheckList

	switch checklist := checklistData.(type) {
	case *types.SimpleCheckList:
		return formatters.FormatSimpleChecklistPreview(checklist)
	case *types.BlockedCheckList:
		return formatters.FormatBlockedChecklistPreview(checklist)
	default:
		return "❌ Неизвестный тип чек-листа"
	}
}
