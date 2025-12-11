package routes

import (
	"log"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/question_edit"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QuestionEditRoutes struct {
	questionEditSvc *question_edit.QuestionEditService
	screenSvc       *screen.ScreenService
}

func NewQuestionEditRoutes(questionEditSvc *question_edit.QuestionEditService, screenSvc *screen.ScreenService) *QuestionEditRoutes {
	return &QuestionEditRoutes{
		questionEditSvc: questionEditSvc,
		screenSvc:       screenSvc,
	}
}

func (r *QuestionEditRoutes) Route(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	log.Printf("[QuestionEditRoutes] 👤 UserID: %d | 💬 Текст: %s | Экран: %s", userID, text, userState.GetCurrentScreen())

	// Проверяем, является ли текст кнопкой с карандашом и номером
	if utils.IsPencilNumberButton(text) {
		r.questionEditSvc.HandleQuestionNumber(userID, update, userState, text)
		return
	}

	switch text {
	case buttons.BtnBack:
		r.questionEditSvc.HandleBack(userID, update, userState)
	case buttons.BtnEditQuestionText:
		r.questionEditSvc.HandleEditQuestionText(userID, update, userState)
	case buttons.BtnEditQuestionType:
		r.questionEditSvc.HandleEditQuestionType(userID, update, userState)
	case buttons.BtnDeleteQuestion:
		r.questionEditSvc.HandleDeleteQuestion(userID, update, userState)
	case buttons.BtnYes:
		r.questionEditSvc.HandleConfirmDelete(userID, update, userState)
	case buttons.BtnNo:
		r.questionEditSvc.HandleCancelDelete(userID, update, userState)
	default:
		// Если текст не кнопка, обрабатываем ввод в зависимости от экрана
		r.handleUserInput(userID, update, userState, text)
	}
}

// handleUserInput обрабатывает текстовый ввод
func (r *QuestionEditRoutes) handleUserInput(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "edit-question-text":
		r.questionEditSvc.HandleNewQuestionTextInput(userID, update, userState, text)
	case "edit-question-type":
		r.questionEditSvc.HandleQuestionTypeSelection(userID, update, userState, text)
	default:
		// Игнорируем ввод на других экранах
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}
