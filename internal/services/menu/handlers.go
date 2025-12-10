package menu

import (
	"context"
	"fmt"
	"log"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleCreateChecklist обрабатывает создание чек-листа
func (r *MenuService) HandleCreateSinpleChecklist(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Очищаем предыдущий чек-лист, если был
	r.stateMgr.ClearCheckList(userID)
	// Переходим к вводу названия
	r.stateMgr.NavigateTo(userID, "create-simple-checklist-name")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleCreateChecklist обрабатывает создание чек-листа
func (r *MenuService) HandleCreateBlocksChecklist(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Очищаем предыдущий чек-лист, если был
	r.stateMgr.ClearCheckList(userID)
	// Переходим к вводу названия
	r.stateMgr.NavigateTo(userID, "create-block-checklist-name")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleMyChecklists обрабатывает просмотр моих чек-листов
func (r *MenuService) HandleMyChecklists(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Получаем черновики пользователя
	ctx := context.Background()
	drafts, err := r.checklistSvc.GetUserDrafts(ctx, userID)
	if err != nil {
		message := fmt.Sprintf("❌ Ошибка при получении чек-листов: %s", err.Error())
		r.screenSvc.SendMessage(update.Message.Chat.ID, message)
		return
	}

	if len(drafts) == 0 {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "📭 У вас пока нет черновиков чек-листов.")
		return
	}

	// Сохраняем список чек-листов в состояние
	userState.Data["my_checklists"] = drafts
	r.stateMgr.SetState(userID, userState)

	// Переходим на экран списка чек-листов
	r.stateMgr.NavigateTo(userID, "my-checklists-list")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[MenuService] ✅ Пользователь %d просмотрел список черновиков (%d шт.)", userID, len(drafts))
}

// HandleLogout обрабатывает выход из системы
func (r *MenuService) HandleLogout(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Логируем начало выхода
	userName := ""
	if userState != nil && userState.User != nil {
		userName = userState.User.FullName
	}
	log.Printf("[MenuRoutes] 🚪 User %d (%s) выходит из системы", userID, userName)

	// Очищаем состояние
	userState.SetUser(nil)
	r.stateMgr.DeleteState(userID)

	// Создаем новое состояние для авторизации
	newState := state.NewUserState(nil, "authorize-admin")
	r.stateMgr.SetState(userID, newState)

	// Отправляем экран авторизации
	r.screenSvc.SendScreen(update.Message.Chat.ID, "authorize-admin", newState)

	log.Printf("[MenuRoutes] ✅ User %d успешно вышел из системы", userID)
}

func (r *MenuService) HandlePublishedChecklists(userID int64, update tgbotapi.Update, userState *state.UserState) {
	r.publishedSvc.HandlePublishedChecklists(userID, update, userState)
}

func (r *MenuService) HandleCanceledChecklists(userID int64, update tgbotapi.Update, userState *state.UserState) {
	r.publishedSvc.HandleUnpublishedChecklists(userID, update, userState)
}
