package answers

import (
	"strings"

	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleAnswerOptionsInput обрабатывает ввод вариантов ответов
func (s *AnswersService) HandleAnswerOptionsInput(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	options := strings.Split(text, "\n")
	cleanOptions := make([]string, 0)

	for _, opt := range options {
		if trimmed := strings.TrimSpace(opt); trimmed != "" {
			cleanOptions = append(cleanOptions, trimmed)
		}
	}

	if len(cleanOptions) < 2 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Нужно минимум 2 варианта")
		return
	}

	err := s.addAnswerOptions(userID, cleanOptions)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ "+err.Error())
		return
	}

	s.stateMgr.NavigateTo(userID, "select-correct-answers")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleCorrectAnswersInput обрабатывает выбор правильных ответов
func (s *AnswersService) HandleCorrectAnswersInput(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	indices, err := s.parseIndices(text)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ "+err.Error())
		return
	}

	// Определяем индекс вопроса
	var idx int
	if lastIdx, err := s.questionSvc.GetLastQuestionIndex(userID); err == nil {
		idx = lastIdx
	} else {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ "+err.Error())
		return
	}

	// Проверяем тип чек-листа и категорию вопроса
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if exists {
		// Определяем категорию вопроса
		var category types.QuestionCategory
		isBlockMode, _ := userState.Data["is_block_mode"].(bool)

		if isBlockMode {
			// Для блоков
			if checklist, ok := checklistData.(*types.BlockedCheckList); ok {
				blockIdx, _ := userState.Data["current_block_index"].(int)
				if blockIdx >= 0 && blockIdx < len(checklist.Blocks) &&
					idx >= 0 && idx < len(checklist.Blocks[blockIdx].Questions) {
					category = checklist.Blocks[blockIdx].Questions[idx].Category
				}
			}
		} else {
			// Для простых чек-листов
			if simpleChecklist, ok := checklistData.(*types.SimpleCheckList); ok && idx >= 0 && idx < len(simpleChecklist.Questions) {
				category = simpleChecklist.Questions[idx].Category
			}
		}

		// Проверяем валидность количества ответов для типа вопроса
		if category == types.CategorySingleChoice && len(indices) != 1 {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Для одиночного выбора нужен 1 правильный ответ")
			return
		}

		if category == types.CategoryMultipleChoice && len(indices) < 2 {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Для множественного выбора нужно минимум 2 правильных ответа")
			return
		}
	}

	err = s.setCorrectAnswers(userID, indices, idx)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ "+err.Error())
		return
	}

	// Завершаем процесс создания/редактирования вопроса
	s.completeQuestionProcess(userID, update, userState)
}

// HandleBack обрабатывает кнопку "Назад" для ответов
func (s *AnswersService) HandleBack(userID int64, update tgbotapi.Update, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "enter-answer-options":
		s.stateMgr.NavigateTo(userID, "enter-question-text")
	case "select-correct-answers":
		s.stateMgr.NavigateTo(userID, "enter-answer-options")
	default:
		// Проверяем режим: блок или простой чек-лист
		isBlockMode, _ := userState.Data["is_block_mode"].(bool)

		if isBlockMode {
			// Возвращаемся в редактор блока
			delete(userState.Data, "is_block_mode")
			delete(userState.Data, "selected_question_type")
			s.stateMgr.SetState(userID, userState)
			s.stateMgr.NavigateTo(userID, "block-editor")
		} else {
			// Проверяем режим: создание или редактирование
			if isEdit, ok := userState.Data["is_edit_mode"].(bool); ok && isEdit {
				// Возвращаемся к деталям вопроса при редактировании
				delete(userState.Data, "is_edit_mode")
				delete(userState.Data, "edit_question_position")
				delete(userState.Data, "selected_question_type")
				s.stateMgr.SetState(userID, userState)
				s.stateMgr.NavigateTo(userID, "edit-question-detail")
			} else {
				// Возвращаемся к выбору типа вопроса при создании
				delete(userState.Data, "selected_question_type")
				s.stateMgr.SetState(userID, userState)
				s.stateMgr.NavigateTo(userID, "select-question-type")
			}
		}
	}

	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleUserInput обрабатывает ввод пользователя для ответов
func (s *AnswersService) HandleUserInput(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "enter-answer-options":
		s.HandleAnswerOptionsInput(userID, update, userState, text)
	case "select-correct-answers":
		s.HandleCorrectAnswersInput(userID, update, userState, text)
	default:
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}
