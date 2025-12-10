package manager

import (
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
	// "time"
)

// StateManager интерфейс менеджера состояний
type StateManager interface {
	// Основные операции
	GetState(userID int64) (*state.UserState, bool)
	SetState(userID int64, state *state.UserState)
	DeleteState(userID int64)

	// Работа с навигацией
	NavigateTo(userID int64, screen string) bool
	GetCurrentScreen(userID int64) (string, bool)

	// Работа с чек-листами
	SetSimpleCheckList(userID int64, checkList *types.SimpleCheckList) bool
	SetBlockedCheckList(userID int64, checkList *types.BlockedCheckList) bool
	GetCheckList(userID int64) (types.CheckListData, bool)
	ClearCheckList(userID int64) bool
	HasCheckList(userID int64) bool
}
