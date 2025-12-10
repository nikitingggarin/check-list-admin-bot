package manager

import "telegram-bot/internal/state_manager/state"

// GetState возвращает состояние пользователя
func (m *MemoryStateManager) GetState(userID int64) (*state.UserState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[userID]
	if !exists {
		return nil, false
	}

	return state, true
}

// SetState устанавливает состояние пользователя
func (m *MemoryStateManager) SetState(userID int64, userState *state.UserState) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[userID] = userState
}

// DeleteState удаляет состояние пользователя
func (m *MemoryStateManager) DeleteState(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, userID)
}

// NavigateTo переводит пользователя на новый экран (ПРОСТО УСТАНАВЛИВАЕМ ЭКРАН)
func (m *MemoryStateManager) NavigateTo(userID int64, screen string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.states[userID]
	if !exists {
		return false
	}

	state.SetCurrentScreen(screen)
	return true
}

// GetCurrentScreen возвращает текущий экран пользователя
func (m *MemoryStateManager) GetCurrentScreen(userID int64) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[userID]
	if !exists {
		return "", false
	}

	return state.GetCurrentScreen(), true
}
