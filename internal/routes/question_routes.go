package routes

import (
	"context"
	"log"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/question"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QuestionRoutes struct {
	questionSvc *question.QuestionService
	screenSvc   *screen.ScreenService
}

func NewQuestionRoutes(questionSvc *question.QuestionService, screenSvc *screen.ScreenService) *QuestionRoutes {
	return &QuestionRoutes{
		questionSvc: questionSvc,
		screenSvc:   screenSvc,
	}
}

func (r *QuestionRoutes) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	log.Printf("[QuestionRoutes] 👤 UserID: %d | 💬 Текст: %s | Экран: %s", userID, text, userState.GetCurrentScreen())

	switch text {
	case buttons.BtnBack:
		r.questionSvc.HandleBack(userID, update, userState)
	case buttons.BtnCompliance:
		r.questionSvc.HandleCompliance(userID, update, userState)
	case buttons.BtnSingleChoice:
		r.questionSvc.HandleSingleChoice(userID, update, userState)
	case buttons.BtnMultipleChoice:
		r.questionSvc.HandleMultipleChoice(userID, update, userState)
	case buttons.BtnTextAnswer:
		r.questionSvc.HandleTextAnswer(userID, update, userState)
	default:
		// Любой другой текст - это ввод текста вопроса
		r.questionSvc.HandleQuestionTextInput(userID, update, userState, text)
	}
}
