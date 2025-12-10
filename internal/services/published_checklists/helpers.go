package published_checklists

import (
	"telegram-bot/internal/state_manager/state"
)

// cleanupPublishedChecklistData очищает данные опубликованного чек-листа
func (s *PublishedChecklistsService) cleanupPublishedChecklistData(userState *state.UserState) {
	delete(userState.Data, "current_published_checklist")
	delete(userState.Data, "current_checklist_type")
	delete(userState.Data, "published_has_blocks")
	delete(userState.Data, "published_total_questions")
	delete(userState.Data, "published_checklist_blocks")
	delete(userState.Data, "published_checklist_questions")
	delete(userState.Data, "published_checklist_answer_options")
}
