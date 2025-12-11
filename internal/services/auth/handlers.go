package auth

import (
	"fmt"
	"log"

	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *AuthService) HandleAuthorization(userID int64, update tgbotapi.Update, userState *state.UserState) error {
	dbUser, err := s.userSvc.AuthorizeAdmin(userID)
	if err != nil {
		return fmt.Errorf("❌ Ошибка авторизации: %s", err)
	}

	stateUser := convertToStateUser(dbUser)
	userState.SetUser(stateUser)
	s.stateMgr.SetState(userID, userState)
	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendScreen(update.Message.Chat.ID, "admin-menu", userState)

	log.Printf("[AuthService] ✅ User %d авторизован как админ", userID)
	return nil
}
