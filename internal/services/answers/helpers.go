package answers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// addAnswerOptions добавляет варианты ответов для нового вопроса
func (s *AnswersService) addAnswerOptions(userID int64, options []string) error {
	idx, err := s.questionSvc.GetLastQuestionIndex(userID)
	if err != nil {
		return err
	}

	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		return fmt.Errorf("чек-лист не найден")
	}

	// Проверяем режим: блок или простой чек-лист
	userState, _ := s.stateMgr.GetState(userID)
	isBlockMode, _ := userState.Data["is_block_mode"].(bool)
	isEditBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockMode || isEditBlockQuestion {
		// ДЛЯ БЛОКОВ (создание или редактирование)
		blockedChecklist, ok := checklistData.(*types.BlockedCheckList)
		if !ok {
			return fmt.Errorf("неверный тип чек-листа для режима блока")
		}

		blockIdx, ok := userState.Data["current_block_index"].(int)
		if !ok {
			return fmt.Errorf("блок не выбран")
		}

		if blockIdx < 0 || blockIdx >= len(blockedChecklist.Blocks) {
			return fmt.Errorf("неверный индекс блока")
		}

		if idx < 0 || idx >= len(blockedChecklist.Blocks[blockIdx].Questions) {
			return fmt.Errorf("неверный индекс вопроса")
		}

		question := &blockedChecklist.Blocks[blockIdx].Questions[idx]

		// Очищаем старые варианты и добавляем новые
		question.AnswerOptions = make([]types.AnswerOption, 0, len(options))
		for _, opt := range options {
			question.AnswerOptions = append(question.AnswerOptions, types.NewAnswerOption(opt, false))
		}

		// Сохраняем количество вариантов
		userState.Data["answer_options_count"] = len(options)
		s.stateMgr.SetState(userID, userState)

		log.Printf("[AnswersService] ✅ Добавлено %d вариантов для вопроса %d в блоке %d",
			len(options), idx+1, blockIdx+1)
	} else {
		// ДЛЯ ПРОСТЫХ ЧЕК-ЛИСТОВ (существующая логика)
		simpleChecklist, ok := checklistData.(*types.SimpleCheckList)
		if !ok {
			return fmt.Errorf("неверный тип чек-листа")
		}

		if idx < 0 || idx >= len(simpleChecklist.Questions) {
			return fmt.Errorf("неверный индекс вопроса")
		}

		question := &simpleChecklist.Questions[idx]

		// Очищаем старые варианты и добавляем новые
		question.AnswerOptions = make([]types.AnswerOption, 0, len(options))
		for _, opt := range options {
			question.AnswerOptions = append(question.AnswerOptions, types.NewAnswerOption(opt, false))
		}

		// Сохраняем количество вариантов
		userState.Data["answer_options_count"] = len(options)
		s.stateMgr.SetState(userID, userState)

		log.Printf("[AnswersService] ✅ Добавлено %d вариантов для вопроса %d", len(options), idx+1)
	}

	return nil
}

// setCorrectAnswers устанавливает правильные ответы
func (s *AnswersService) setCorrectAnswers(userID int64, indices []int, questionIdx int) error {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		return fmt.Errorf("чек-лист не найден")
	}

	// Проверяем режим: блок или простой чек-лист
	userState, _ := s.stateMgr.GetState(userID)
	isBlockMode, _ := userState.Data["is_block_mode"].(bool)
	isEditBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockMode || isEditBlockQuestion {
		// ДЛЯ БЛОКОВ
		blockedChecklist, ok := checklistData.(*types.BlockedCheckList)
		if !ok {
			return fmt.Errorf("неверный тип чек-листа для режима блока")
		}

		blockIdx, ok := userState.Data["current_block_index"].(int)
		if !ok {
			return fmt.Errorf("блок не выбран")
		}

		if blockIdx < 0 || blockIdx >= len(blockedChecklist.Blocks) {
			return fmt.Errorf("неверный индекс блока")
		}

		if questionIdx < 0 || questionIdx >= len(blockedChecklist.Blocks[blockIdx].Questions) {
			return fmt.Errorf("неверный индекс вопроса")
		}

		question := &blockedChecklist.Blocks[blockIdx].Questions[questionIdx]

		// Получаем количество вариантов
		count := len(question.AnswerOptions)

		// Проверка индексов
		for _, i := range indices {
			if i < 0 || i >= count {
				return fmt.Errorf("неверный индекс: %d", i+1)
			}
		}

		// Установка правильных ответов
		for i := range question.AnswerOptions {
			question.AnswerOptions[i].IsCorrect = false
		}
		for _, correctIdx := range indices {
			if correctIdx < len(question.AnswerOptions) {
				question.AnswerOptions[correctIdx].IsCorrect = true
			}
		}

		log.Printf("[AnswersService] ✅ Установлены правильные ответы: %v для вопроса %d в блоке %d",
			indices, questionIdx+1, blockIdx+1)
	} else {
		// ДЛЯ ПРОСТЫХ ЧЕК-ЛИСТОВ (существующая логика)
		simpleChecklist, ok := checklistData.(*types.SimpleCheckList)
		if !ok {
			return fmt.Errorf("неверный тип чек-листа")
		}

		if questionIdx < 0 || questionIdx >= len(simpleChecklist.Questions) {
			return fmt.Errorf("неверный индекс вопроса")
		}

		question := &simpleChecklist.Questions[questionIdx]

		// Получаем количество вариантов
		count := len(question.AnswerOptions)

		// Проверка индексов
		for _, i := range indices {
			if i < 0 || i >= count {
				return fmt.Errorf("неверный индекс: %d", i+1)
			}
		}

		// Установка правильных ответов
		for i := range question.AnswerOptions {
			question.AnswerOptions[i].IsCorrect = false
		}
		for _, correctIdx := range indices {
			if correctIdx < len(question.AnswerOptions) {
				question.AnswerOptions[correctIdx].IsCorrect = true
			}
		}

		log.Printf("[AnswersService] ✅ Установлены правильные ответы: %v для вопроса %d",
			indices, questionIdx+1)
	}

	return nil
}

// completeQuestionProcess завершает процесс создания/редактирования вопроса
func (s *AnswersService) completeQuestionProcess(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Проверяем режим: блок или простой чек-лист
	isBlockMode, _ := userState.Data["is_block_mode"].(bool)
	isEditBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	// Очищаем временные данные
	s.cleanupQuestionData(userID)

	if isBlockMode || isEditBlockQuestion {
		// РЕЖИМ БЛОКА
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
			log.Printf("[AnswersService] ✅ Редактирование вопроса в блоке завершено")
		} else {
			// СОЗДАНИЕ вопроса в блоке
			// НЕ удаляем is_block_mode - остаемся в режиме добавления
			delete(userState.Data, "selected_question_type")
			s.stateMgr.SetState(userID, userState)

			// Возвращаемся к выбору типа вопроса для добавления следующего
			s.stateMgr.NavigateTo(userID, "select-question-type")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			log.Printf("[AnswersService] ✅ Создание вопроса в блоке завершено. Возвращаемся в select-question-type")
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
			log.Printf("[AnswersService] ✅ Редактирование вопроса завершено")
		} else {
			// Очищаем данные создания
			delete(userState.Data, "selected_question_type")
			s.stateMgr.SetState(userID, userState)

			// Возвращаемся к выбору типа вопроса для создания следующего
			s.stateMgr.NavigateTo(userID, "select-question-type")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			log.Printf("[AnswersService] ✅ Создание вопроса завершено")
		}
	}
}

// parseIndices парсит номера ответов
func (s *AnswersService) parseIndices(text string) ([]int, error) {
	parts := strings.Split(text, ",")
	indices := make([]int, 0)
	seen := make(map[int]bool)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("неверный номер: '%s'", part)
		}

		if num < 1 {
			return nil, fmt.Errorf("номер должен быть положительным")
		}

		idx := num - 1 // Конвертируем в 0-based

		if !seen[idx] {
			seen[idx] = true
			indices = append(indices, idx)
		}
	}

	if len(indices) == 0 {
		return nil, fmt.Errorf("укажите номера правильных ответов (например: 1 или 1,3)")
	}

	return indices, nil
}

// cleanupQuestionData очищает временные данные
func (s *AnswersService) cleanupQuestionData(userID int64) {
	userState, exists := s.stateMgr.GetState(userID)
	if !exists {
		return
	}

	delete(userState.Data, "last_question_index")
	delete(userState.Data, "answer_options_count")
	s.stateMgr.SetState(userID, userState)
}
