package routes

import (
	"context"
	"log"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/auth"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AuthRoutes struct {
	authSvc   *auth.AuthService
	screenSvc *screen.ScreenService
}

func NewAuthRoutes(authSvc *auth.AuthService, screenSvc *screen.ScreenService) *AuthRoutes {
	return &AuthRoutes{
		authSvc:   authSvc,
		screenSvc: screenSvc,
	}
}

func (r *AuthRoutes) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	log.Printf("[AuthRoutes] 👤 UserID: %d | 💬 Текст: %s", userID, text)

	switch text {
	case buttons.BtnAuth:
		r.authSvc.HandleAuthorization(ctx, userID, update, userState)

	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}
