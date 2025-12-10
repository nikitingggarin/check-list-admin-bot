package routes

import (
	"context"
	"log"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/my_checklists"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MyChecklistsRoutes struct {
	myChecklistsSvc *my_checklists.MyChecklistsService
	screenSvc       *screen.ScreenService
}

func NewMyChecklistsRoutes(myChecklistsSvc *my_checklists.MyChecklistsService, screenSvc *screen.ScreenService) *MyChecklistsRoutes {
	return &MyChecklistsRoutes{
		myChecklistsSvc: myChecklistsSvc,
		screenSvc:       screenSvc,
	}
}

func (r *MyChecklistsRoutes) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	log.Printf("[MyChecklistsRoutes] 👤 UserID: %d | 💬 Текст: %s | Экран: %s", userID, text, userState.GetCurrentScreen())

	// Проверяем, является ли текст кнопкой с карандашом и номером (формат: "✏️ 1")
	if utils.IsPencilNumberButton(text) {
		r.myChecklistsSvc.HandleChecklistNumber(userID, update, userState, text)
		return
	}

	switch text {
	case buttons.BtnBack:
		r.handleBack(userID, update, userState)
	case buttons.BtnEditCheckList:
		r.myChecklistsSvc.HandleEditChecklist(userID, update, userState)
	case buttons.BtnDeleteCheckList:
		r.myChecklistsSvc.HandleDeleteChecklist(userID, update, userState)
	case buttons.BtnPublishChecklist:
		r.myChecklistsSvc.HandlePublishChecklist(userID, update, userState)
	case buttons.BtnYes:
		r.myChecklistsSvc.HandleConfirmDelete(userID, update, userState)
	case buttons.BtnNo:
		r.myChecklistsSvc.HandleCancelDelete(userID, update, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

func (r *MyChecklistsRoutes) handleBack(userID int64, update tgbotapi.Update, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "my-checklists-list":
		r.myChecklistsSvc.HandleBackFromList(userID, update, userState)
	case "checklist-detail":
		// Возвращаемся к списку чек-листов
		r.myChecklistsSvc.HandleBackFromDetail(userID, update, userState)
	case "confirm-delete-checklist":
		// Возвращаемся к деталям чек-листа
		r.screenSvc.SendScreen(update.Message.Chat.ID, "checklist-detail", userState)
	default:
		r.screenSvc.SendScreen(update.Message.Chat.ID, "my-checklists-list", userState)
	}
}
