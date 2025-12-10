package question

import (
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleBack обрабатывает кнопку "Назад" для вопросов
func (s *QuestionService) HandleBack(userID int64, update tgbotapi.Update, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "select-question-type":
		// Определяем куда возвращаться в зависимости от режима
		isBlockMode, _ := userState.Data["is_block_mode"].(bool)
		if isBlockMode {
			// Возвращаемся в редактор блока
			delete(userState.Data, "is_block_mode")
			s.stateMgr.SetState(userID, userState)
			s.stateMgr.NavigateTo(userID, "block-editor")
		} else if isEdit, ok := userState.Data["is_edit_mode"].(bool); ok && isEdit {
			// Редактирование вопроса в простом чек-листе
			delete(userState.Data, "is_edit_mode")
			delete(userState.Data, "edit_question_position")
			s.stateMgr.SetState(userID, userState)
			s.stateMgr.NavigateTo(userID, "edit-question-detail")
		} else {
			// Создание вопроса в простом чек-листе
			s.stateMgr.NavigateTo(userID, "simple-checklist-editor")
		}
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	case "enter-question-text":
		// Возвращаемся к выбору типа вопроса
		s.stateMgr.NavigateTo(userID, "select-question-type")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	default:
		s.stateMgr.NavigateTo(userID, "simple-checklist-editor")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// HandleCompliance обрабатывает выбор типа "Соответствие"
func (s *QuestionService) HandleCompliance(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.handleQuestionType(userID, update, userState, types.CategoryCompliance)
}

// HandleSingleChoice обрабатывает выбор типа "Одиночный выбор"
func (s *QuestionService) HandleSingleChoice(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.handleQuestionType(userID, update, userState, types.CategorySingleChoice)
}

// HandleMultipleChoice обрабатывает выбор типа "Множественный выбор"
func (s *QuestionService) HandleMultipleChoice(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.handleQuestionType(userID, update, userState, types.CategoryMultipleChoice)
}

// HandleTextAnswer обрабатывает выбор типа "Текстовый ответ"
func (s *QuestionService) HandleTextAnswer(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.handleQuestionType(userID, update, userState, types.CategoryTextAnswer)
}

// HandleQuestionTextInput обрабатывает ввод текста вопроса
func (s *QuestionService) HandleQuestionTextInput(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	categoryStr, ok := userState.Data["selected_question_type"].(string)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Тип вопроса не найден")
		s.stateMgr.NavigateTo(userID, "select-question-type")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
		return
	}

	err := s.createQuestion(userID, text, types.QuestionCategory(categoryStr))
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ "+err.Error())
		return
	}

	category := types.QuestionCategory(categoryStr)

	// Только для SingleChoice и MultipleChoice нужны варианты ответов
	if category == types.CategorySingleChoice || category == types.CategoryMultipleChoice {
		s.stateMgr.NavigateTo(userID, "enter-answer-options")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	} else {
		// Для Compliance и TextAnswer - завершаем процесс
		s.completeQuestionCreation(userID, update, userState)
	}
}
