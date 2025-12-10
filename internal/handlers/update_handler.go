package handlers

import (
	"context"
	"fmt"
	"log"

	"telegram-bot/internal/state_manager/debug"
	"telegram-bot/internal/state_manager/manager"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UpdateHandler обрабатывает входящие обновления Telegram
type UpdateHandler struct {
	stateMgr manager.StateManager
	router   interface {
		Route(ctx context.Context, userID int64, update tgbotapi.Update, text string)
	}
}

// NewUpdateHandler создает новый обработчик обновлений
func NewUpdateHandler(router interface {
	Route(ctx context.Context, userID int64, update tgbotapi.Update, text string)
}, stateMgr manager.StateManager) *UpdateHandler {
	return &UpdateHandler{
		stateMgr: stateMgr,
		router:   router,
	}
}

// HandleUpdate обрабатывает одно обновление Telegram
func (h *UpdateHandler) HandleUpdate(update tgbotapi.Update) {
	// Пропускаем не текстовые сообщения
	if update.Message == nil {
		return
	}
	ctx := context.Background()
	userID := update.Message.From.ID
	text := update.Message.Text
	// Передаем в роутер
	h.router.Route(ctx, userID, update, text)

	// Получаем и логируем состояние ПОСЛЕ обработки
	log.Println("\n📊 ТЕКУЩЕЕ СОСТОЯНИЕ:")
	if state, exists := h.stateMgr.GetState(userID); exists {
		fmt.Print(state)
		debug.PrintState(state)
	}

}
