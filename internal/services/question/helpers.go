package question

import (
	"fmt"
	"log"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleQuestionType общая логика выбора типа
func (s *QuestionService) handleQuestionType(userID int64, update tgbotapi.Update, userState *state.UserState, category types.QuestionCategory) {
	// Всегда перезаписываем тип
	userState.Data["selected_question_type"] = string(category)
	s.stateMgr.SetState(userID, userState)
	s.stateMgr.NavigateTo(userID, "enter-question-text")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// createQuestion создает или заменяет вопрос
func (s *QuestionService) createQuestion(userID int64, text string, category types.QuestionCategory) error {
	userState, exists := s.stateMgr.GetState(userID)
	if !exists {
		return fmt.Errorf("состояние не найдено")
	}

	// Проверяем режим: блок или простой чек-лист
	isBlockMode, _ := userState.Data["is_block_mode"].(bool)
	isEditBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockMode || isEditBlockQuestion {
		// РЕЖИМ БЛОКА: добавляем вопрос в текущий блок
		return s.createQuestionInBlock(userID, text, category, userState, isEditBlockQuestion)
	} else {
		// РЕЖИМ ПРОСТОГО ЧЕК-ЛИСТА: существующая логика
		return s.createQuestionInSimpleChecklist(userID, text, category, userState)
	}
}

// createQuestionInBlock добавляет вопрос в текущий блок
func (s *QuestionService) createQuestionInBlock(userID int64, text string, category types.QuestionCategory, userState *state.UserState, isEditMode bool) error {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		return fmt.Errorf("чек-лист не найден")
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		return fmt.Errorf("неверный тип чек-листа для режима блока")
	}

	// Получаем индекс текущего блока
	blockIndex, ok := userState.Data["current_block_index"].(int)
	if !ok {
		return fmt.Errorf("блок не выбран")
	}

	if blockIndex < 0 || blockIndex >= len(checklist.Blocks) {
		return fmt.Errorf("неверный индекс блока")
	}

	// Проверяем, редактируем ли существующий вопрос в блоке
	if isEditMode {
		// РЕДАКТИРОВАНИЕ вопроса в блоке
		if position, ok := userState.Data["edit_question_position"].(int); ok {
			if position >= 0 && position < len(checklist.Blocks[blockIndex].Questions) {
				// Создаем новый вопрос
				question := types.NewQuestion(text, category)

				// Заменяем вопрос на позиции
				checklist.Blocks[blockIndex].Questions[position] = question

				// Сохраняем индекс последнего созданного вопроса для работы с ответами
				userState.Data["last_question_index"] = position
				s.stateMgr.SetState(userID, userState)

				log.Printf("[QuestionService] ✅ Вопрос %d в блоке %d заменен: '%s' (%s)",
					position+1, blockIndex+1, text, category)
				return nil
			}
		}
	}

	// СОЗДАНИЕ нового вопроса в блоке
	question := types.NewQuestion(text, category)

	// Добавляем вопрос в блок
	checklist.Blocks[blockIndex].AddQuestion(question)

	// Сохраняем индекс последнего созданного вопроса
	lastQuestionIndex := len(checklist.Blocks[blockIndex].Questions) - 1
	userState.Data["last_question_index"] = lastQuestionIndex
	s.stateMgr.SetState(userID, userState)

	log.Printf("[QuestionService] ✅ Вопрос добавлен в блок %d: '%s' (%s)",
		blockIndex+1, text, category)
	return nil
}

// createQuestionInSimpleChecklist добавляет вопрос в простой чек-лист
func (s *QuestionService) createQuestionInSimpleChecklist(userID int64, text string, category types.QuestionCategory, userState *state.UserState) error {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		return fmt.Errorf("чек-лист не найден")
	}

	simpleChecklist, ok := checklistData.(*types.SimpleCheckList)
	if !ok {
		return fmt.Errorf("неверный тип чек-листа")
	}

	// Проверяем, редактируем ли существующий вопрос
	if isEdit, ok := userState.Data["is_edit_mode"].(bool); ok && isEdit {
		// РЕДАКТИРОВАНИЕ: Заменяем вопрос на заданной позиции
		if position, ok := userState.Data["edit_question_position"].(int); ok {
			if position >= 0 && position < len(simpleChecklist.Questions) {
				// Создаем новый вопрос
				question := types.NewQuestion(text, category)

				// Заменяем вопрос на позиции
				simpleChecklist.Questions[position] = question

				// Сохраняем индекс последнего созданного вопроса для работы с ответами
				userState.Data["last_question_index"] = position
				s.stateMgr.SetState(userID, userState)

				log.Printf("[QuestionService] ✅ Вопрос %d заменен: '%s' (%s)", position+1, text, category)
				return nil
			}
		}
	}

	// СОЗДАНИЕ: Добавляем новый вопрос в конец
	question := types.NewQuestion(text, category)
	simpleChecklist.AddQuestion(question)
	userState.Data["last_question_index"] = len(simpleChecklist.Questions) - 1
	s.stateMgr.SetState(userID, userState)

	log.Printf("[QuestionService] ✅ Вопрос создан: '%s' (%s)", text, category)
	return nil
}

// completeQuestionCreation завершает процесс создания/редактирования вопроса
func (s *QuestionService) completeQuestionCreation(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Проверяем режим: блок или простой чек-лист
	isBlockMode, _ := userState.Data["is_block_mode"].(bool)
	isEditBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockMode || isEditBlockQuestion {
		// РЕЖИМ БЛОКА (создание или редактирование вопроса)
		if isEditBlockQuestion {
			// Редактирование вопроса в блоке
			delete(userState.Data, "is_edit_block_questions")
			delete(userState.Data, "is_edit_mode")
			delete(userState.Data, "edit_question_position")
			delete(userState.Data, "selected_question_type")
			s.stateMgr.SetState(userID, userState)

			// Возвращаемся к списку вопросов блока
			s.stateMgr.NavigateTo(userID, "block-view-questions")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			log.Printf("[QuestionService] ✅ Редактирование вопроса в блоке завершено")
		} else {
			// СОЗДАНИЕ вопроса в блоке
			// НЕ удаляем is_block_mode - остаемся в режиме добавления
			delete(userState.Data, "selected_question_type")
			s.stateMgr.SetState(userID, userState)

			// Возвращаемся к выбору типа вопроса для добавления следующего
			s.stateMgr.NavigateTo(userID, "select-question-type")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			log.Printf("[QuestionService] ✅ Создание вопроса в блоке завершено. Возвращаемся в select-question-type")
		}
	} else {
		// РЕЖИМ ПРОСТОГО ЧЕК-ЛИСТА
		// Проверяем режим: создание или редактирование
		if isEdit, ok := userState.Data["is_edit_mode"].(bool); ok && isEdit {
			// Очищаем данные редактирования
			delete(userState.Data, "is_edit_mode")
			delete(userState.Data, "edit_question_position")
			delete(userState.Data, "selected_question_type")
			s.stateMgr.SetState(userID, userState)

			// Возвращаемся к списку вопросов после редактирования
			s.stateMgr.NavigateTo(userID, "view-question")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			log.Printf("[QuestionService] ✅ Редактирование вопроса завершено")
		} else {
			// Очищаем данные создания
			delete(userState.Data, "selected_question_type")
			s.stateMgr.SetState(userID, userState)

			// Возвращаемся к выбору типа вопроса для создания следующего
			s.stateMgr.NavigateTo(userID, "select-question-type")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			log.Printf("[QuestionService] ✅ Создание вопроса завершено")
		}
	}
}

// GetLastQuestionIndex возвращает индекс последнего вопроса
func (s *QuestionService) GetLastQuestionIndex(userID int64) (int, error) {
	userState, exists := s.stateMgr.GetState(userID)
	if !exists {
		return -1, fmt.Errorf("состояние не найдено")
	}

	idx, ok := userState.Data["last_question_index"].(int)
	if !ok {
		return -1, fmt.Errorf("индекс вопроса не найден")
	}

	return idx, nil
}
