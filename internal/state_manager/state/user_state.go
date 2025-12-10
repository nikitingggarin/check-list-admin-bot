package state

import (
	"telegram-bot/internal/state_manager/types"
)

// UserState - полное состояние пользователя
type UserState struct {
	User             *types.User            `json:"user"`
	CurrentScreen    string                 `json:"current_screen"`
	CurrentCheckList types.CheckListData    `json:"current_check_list,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
}

// NewUserState создает новое состояние пользователя
func NewUserState(user *types.User, initialScreen string) *UserState {
	if user == nil {
		user = &types.User{
			ID:         -1,
			TelegramID: -1,
			Username:   "",
			FullName:   "",
			Role:       types.RoleUser,
		}
	}

	state := &UserState{
		User:          user,
		CurrentScreen: initialScreen,
		Data:          make(map[string]interface{}),
	}

	return state
}

// SetUser обновляет пользователя
func (s *UserState) SetUser(user *types.User) {
	s.User = user
}

// SetCurrentScreen устанавливает текущий экран
func (s *UserState) SetCurrentScreen(screen string) {
	s.CurrentScreen = screen
}

// GetCurrentScreen возвращает текущий экран
func (s *UserState) GetCurrentScreen() string {
	return s.CurrentScreen
}

// ClearCheckList очищает текущий чек-лист
func (s *UserState) ClearCheckList() {
	s.CurrentCheckList = nil
}

// HasCheckList проверяет наличие чек-листа
func (s *UserState) HasCheckList() bool {
	return s.CurrentCheckList != nil
}
