package manager

import "telegram-bot/internal/state_manager/types"

// SetSimpleCheckList устанавливает простой чек-лист
func (m *MemoryStateManager) SetSimpleCheckList(userID int64, checkList *types.SimpleCheckList) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.states[userID]
	if !exists {
		return false
	}

	state.SetSimpleCheckList(checkList)
	return true
}

// SetBlockedCheckList устанавливает чек-лист с блоками
func (m *MemoryStateManager) SetBlockedCheckList(userID int64, checkList *types.BlockedCheckList) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.states[userID]
	if !exists {
		return false
	}

	state.SetBlockedCheckList(checkList)
	return true
}

// GetCheckList возвращает текущий чек-лист
func (m *MemoryStateManager) GetCheckList(userID int64) (types.CheckListData, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[userID]
	if !exists || state.CurrentCheckList == nil {
		return nil, false
	}

	return state.CurrentCheckList, true
}

// ClearCheckList очищает текущий чек-лист
func (m *MemoryStateManager) ClearCheckList(userID int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.states[userID]
	if !exists {
		return false
	}

	state.ClearCheckList()
	return true
}

// HasCheckList проверяет наличие чек-листа
func (m *MemoryStateManager) HasCheckList(userID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state, exists := m.states[userID]
	if !exists {
		return false
	}

	return state.HasCheckList()
}
