package state

import (
	"telegram-bot/internal/state_manager/types"
)

// SetSimpleCheckList устанавливает простой чек-лист
func (s *UserState) SetSimpleCheckList(checkList *types.SimpleCheckList) {
	s.CurrentCheckList = checkList
}

// SetBlockedCheckList устанавливает чек-лист с блоками
func (s *UserState) SetBlockedCheckList(checkList *types.BlockedCheckList) {
	s.CurrentCheckList = checkList
}
