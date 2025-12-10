package manager

import (
	"sync"
	"telegram-bot/internal/state_manager/state"
)

// MemoryStateManager - реализация StateManager
type MemoryStateManager struct {
	states    map[int64]*state.UserState
	mu        sync.RWMutex
	stopChan  chan struct{}
	isRunning bool
}

// NewMemoryStateManager создает новый менеджер состояний
func NewMemoryStateManager() *MemoryStateManager {
	return &MemoryStateManager{
		states:   make(map[int64]*state.UserState),
		stopChan: make(chan struct{}),
	}
}
