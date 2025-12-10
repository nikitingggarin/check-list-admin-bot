package routes

import (
	"context"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/menu"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MenuRoutes содержит методы для маршрутизации главного меню
type MenuRoutes struct {
	menuSvc   *menu.MenuService
	screenSvc *screen.ScreenService
}

// NewMenuRoutes создает новый роутер для главного меню
func NewMenuRoutes(menuSvc *menu.MenuService, screenSvc *screen.ScreenService) *MenuRoutes {
	return &MenuRoutes{
		menuSvc:   menuSvc,
		screenSvc: screenSvc,
	}
}

// Route маршрутизирует сообщение в главном меню
func (r *MenuRoutes) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string, userState *state.UserState) {

	switch text {
	case buttons.BtnCreateSimpleChecklist:
		r.menuSvc.HandleCreateSinpleChecklist(userID, update, userState)
	case buttons.BtnCreateBlockChecklist:
		r.menuSvc.HandleCreateBlocksChecklist(userID, update, userState)
	case buttons.BtnMyChecklists:
		r.menuSvc.HandleMyChecklists(userID, update, userState)
	case buttons.BtnPublished:
		r.menuSvc.HandlePublishedChecklists(userID, update, userState)
	case buttons.BtnCanceled:
		r.menuSvc.HandleCanceledChecklists(userID, update, userState)
	case buttons.BtnLogout:
		r.menuSvc.HandleLogout(userID, update, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}
