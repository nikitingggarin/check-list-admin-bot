package routes

import (
	"log"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/answers"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AnswersRoutes struct {
	answersSvc *answers.AnswersService
	screenSvc  *screen.ScreenService
}

func NewAnswersRoutes(answersSvc *answers.AnswersService, screenSvc *screen.ScreenService) *AnswersRoutes {
	return &AnswersRoutes{
		answersSvc: answersSvc,
		screenSvc:  screenSvc,
	}
}

func (r *AnswersRoutes) Route(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	log.Printf("[AnswersRoutes] 👤 UserID: %d | 💬 Текст: %s | Экран: %s", userID, text, userState.GetCurrentScreen())

	switch text {
	case buttons.BtnBack:
		r.answersSvc.HandleBack(userID, update, userState)
	default:
		// Любой другой текст - это ввод вариантов или правильных ответов
		r.answersSvc.HandleUserInput(userID, update, userState, text)
	}
}
