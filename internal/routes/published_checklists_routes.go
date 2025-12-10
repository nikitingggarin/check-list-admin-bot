package routes

import (
	"context"
	"log"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/published_checklists"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PublishedChecklistsRoutes struct {
	publishedSvc *published_checklists.PublishedChecklistsService
	screenSvc    *screen.ScreenService
}

func NewPublishedChecklistsRoutes(publishedSvc *published_checklists.PublishedChecklistsService, screenSvc *screen.ScreenService) *PublishedChecklistsRoutes {
	return &PublishedChecklistsRoutes{
		publishedSvc: publishedSvc,
		screenSvc:    screenSvc,
	}
}

func (r *PublishedChecklistsRoutes) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	log.Printf("[PublishedChecklistsRoutes] 👤 UserID: %d | 💬 Текст: %s | Экран: %s", userID, text, userState.GetCurrentScreen())

	// Проверяем, является ли текст кнопкой с карандашом и номером
	if utils.IsPencilNumberButton(text) {
		r.publishedSvc.HandleChecklistNumber(userID, update, userState, text)
		return
	}

	switch text {
	case buttons.BtnBack:
		r.handleBack(userID, update, userState)
	case buttons.BtnUnPublish:
		r.publishedSvc.HandleUnpublish(userID, update, userState)
	case buttons.BtnPublish:
		r.publishedSvc.HandleRepublish(userID, update, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

func (r *PublishedChecklistsRoutes) handleBack(userID int64, update tgbotapi.Update, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "published-checklists-list":
		r.publishedSvc.HandleBackFromList(userID, update, userState)
	case "published-checklist-detail":
		r.publishedSvc.HandleBackFromDetail(userID, update, userState)
	default:
		r.screenSvc.SendScreen(update.Message.Chat.ID, "published-checklists-list", userState)
	}
}
