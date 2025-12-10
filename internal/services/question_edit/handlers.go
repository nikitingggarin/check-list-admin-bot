package question_edit

import (
	"log"
	"strconv"
	"strings"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleEditQuestions начинает редактирование вопросов (показывает список)
func (s *QuestionEditService) HandleEditQuestions(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Проверяем, редактируем ли вопросы в блоке
	isBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockQuestion {
		// Редактирование вопросов в блоке
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		blockedChecklist := checklistData.(*types.BlockedCheckList)
		blockIdx, _ := userState.Data["current_block_index"].(int)

		// Сохраняем индекс редактируемого вопроса (первый по умолчанию)
		userState.Data["edit_question_index"] = 0
		userState.Data["total_questions"] = len(blockedChecklist.Blocks[blockIdx].Questions)
		s.stateMgr.SetState(userID, userState)

		// Переходим к просмотру списка вопросов блока
		s.stateMgr.NavigateTo(userID, "block-view-questions")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	} else {
		// Редактирование вопросов в простом чек-листе (существующая логика)
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		simpleChecklist := checklistData.(*types.SimpleCheckList)

		// Сохраняем индекс редактируемого вопроса (первый по умолчанию)
		userState.Data["edit_question_index"] = 0
		userState.Data["total_questions"] = len(simpleChecklist.Questions)
		s.stateMgr.SetState(userID, userState)

		// Переходим к просмотру списка вопросов
		s.stateMgr.NavigateTo(userID, "view-question")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// HandleBack обрабатывает кнопку "Назад" при редактировании вопросов
func (s *QuestionEditService) HandleBack(userID int64, update tgbotapi.Update, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "view-question":
		// Возвращаемся в редактор простого чек-листа
		s.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	case "block-view-questions":
		// Возвращаемся в редактор блока
		delete(userState.Data, "is_edit_block_questions")
		s.stateMgr.SetState(userID, userState)
		s.stateMgr.NavigateTo(userID, "block-editor")
	case "edit-question-detail":
		// Возвращаемся к списку вопросов
		isBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)
		if isBlockQuestion {
			s.stateMgr.NavigateTo(userID, "block-view-questions")
		} else {
			s.stateMgr.NavigateTo(userID, "view-question")
		}
	case "edit-question-text", "edit-question-type", "confirm-delete-question":
		// Возвращаемся к деталям вопроса
		s.stateMgr.NavigateTo(userID, "edit-question-detail")
	default:
		s.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	}

	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleQuestionNumber обрабатывает выбор номера вопроса
func (s *QuestionEditService) HandleQuestionNumber(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	// Извлекаем номер из кнопки с карандашом
	cleanText := strings.TrimPrefix(text, "✏️")
	cleanText = strings.TrimSpace(cleanText)

	number, err := strconv.Atoi(cleanText)
	if err != nil || number < 1 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный номер вопроса")
		return
	}

	// Проверяем, редактируем ли вопросы в блоке
	isBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockQuestion {
		// Редактирование вопроса в блоке
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		blockedChecklist := checklistData.(*types.BlockedCheckList)
		blockIdx, _ := userState.Data["current_block_index"].(int)

		if number > len(blockedChecklist.Blocks[blockIdx].Questions) {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Вопрос с таким номером не существует")
			return
		}
	} else {
		// Редактирование вопроса в простом чек-листе
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		simpleChecklist := checklistData.(*types.SimpleCheckList)

		if number > len(simpleChecklist.Questions) {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Вопрос с таким номером не существует")
			return
		}
	}

	// Сохраняем индекс выбранного вопроса
	userState.Data["edit_question_index"] = number - 1
	s.stateMgr.SetState(userID, userState)

	// Переходим на экран редактирования конкретного вопроса
	s.stateMgr.NavigateTo(userID, "edit-question-detail")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[QuestionEditService] ✅ Выбран вопрос %d для редактирования", number)
}

// HandleEditQuestionText начинает редактирование текста вопроса
func (s *QuestionEditService) HandleEditQuestionText(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Переходим к редактированию текста
	s.stateMgr.NavigateTo(userID, "edit-question-text")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleEditQuestionType начинает процесс изменения типа вопроса
func (s *QuestionEditService) HandleEditQuestionType(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Получаем индекс редактируемого вопроса
	idx, ok := userState.Data["edit_question_index"].(int)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка: не найден индекс вопроса")
		return
	}

	// Сохраняем индекс и позицию вопроса для редактирования
	userState.Data["edit_question_position"] = idx
	userState.Data["is_edit_mode"] = true
	s.stateMgr.SetState(userID, userState)

	// Переходим к выбору типа вопроса (полный флоу создания с нуля)
	s.stateMgr.NavigateTo(userID, "select-question-type")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[QuestionEditService] 🚀 Начинаем изменение типа вопроса %d", idx+1)
}

// HandleDeleteQuestion начинает процесс удаления вопроса
func (s *QuestionEditService) HandleDeleteQuestion(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Переходим к подтверждению удаления
	s.stateMgr.NavigateTo(userID, "confirm-delete-question")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleNewQuestionTextInput обрабатывает ввод нового текста вопроса
func (s *QuestionEditService) HandleNewQuestionTextInput(userID int64, update tgbotapi.Update, userState *state.UserState, newText string) {
	// Проверяем минимальную длину
	if len(strings.TrimSpace(newText)) < 3 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Текст вопроса должен содержать хотя бы 3 символа")
		return
	}

	// Получаем индекс вопроса
	idx, ok := userState.Data["edit_question_index"].(int)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка: не найден индекс вопроса")
		return
	}

	// Проверяем, редактируем ли вопрос в блоке
	isBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockQuestion {
		// Редактирование вопроса в блоке
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		checklist := checklistData.(*types.BlockedCheckList)
		blockIdx, _ := userState.Data["current_block_index"].(int)

		// Проверяем диапазон
		if idx < 0 || idx >= len(checklist.Blocks[blockIdx].Questions) {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка: неверный индекс вопроса")
			return
		}

		// Обновляем текст вопроса
		checklist.Blocks[blockIdx].Questions[idx].Text = newText

		// Возвращаемся к деталям вопроса
		s.stateMgr.NavigateTo(userID, "edit-question-detail")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

		log.Printf("[QuestionEditService] ✅ Текст вопроса %d в блоке изменен", idx+1)
	} else {
		// Редактирование вопроса в простом чек-листе
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		simpleChecklist := checklistData.(*types.SimpleCheckList)

		// Проверяем диапазон
		if idx < 0 || idx >= len(simpleChecklist.Questions) {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка: неверный индекс вопроса")
			return
		}

		// Обновляем текст вопроса
		simpleChecklist.Questions[idx].Text = newText

		// Возвращаемся к деталям вопроса
		s.stateMgr.NavigateTo(userID, "edit-question-detail")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

		log.Printf("[QuestionEditService] ✅ Текст вопроса %d изменен", idx+1)
	}
}

// HandleQuestionTypeSelection обрабатывает выбор типа вопроса (для создания/редактирования)
func (s *QuestionEditService) HandleQuestionTypeSelection(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	var category types.QuestionCategory

	switch text {
	case buttons.BtnCompliance:
		category = types.CategoryCompliance
	case buttons.BtnSingleChoice:
		category = types.CategorySingleChoice
	case buttons.BtnMultipleChoice:
		category = types.CategoryMultipleChoice
	case buttons.BtnTextAnswer:
		category = types.CategoryTextAnswer
	case buttons.BtnBack:
		// Проверяем режим: создание или редактирование
		if isEdit, ok := userState.Data["is_edit_mode"].(bool); ok && isEdit {
			// Возвращаемся к деталям вопроса при редактировании
			delete(userState.Data, "is_edit_mode")
			delete(userState.Data, "edit_question_position")
			s.stateMgr.SetState(userID, userState)
			s.stateMgr.NavigateTo(userID, "edit-question-detail")
		} else {
			// Возвращаемся в редактор при создании
			s.stateMgr.NavigateTo(userID, "simple-checklist-editor")
		}
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
		return
	default:
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
		return
	}

	// Сохраняем выбранный тип вопроса
	userState.Data["selected_question_type"] = string(category)
	s.stateMgr.SetState(userID, userState)

	// Переходим к вводу текста вопроса
	s.stateMgr.NavigateTo(userID, "enter-question-text")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[QuestionEditService] ✅ Выбран тип вопроса: %s", category)
}

// HandleConfirmDelete подтверждает удаление вопроса
func (s *QuestionEditService) HandleConfirmDelete(userID int64, update tgbotapi.Update, userState *state.UserState) {
	idx, ok := userState.Data["edit_question_index"].(int)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка: не найден индекс вопроса")
		return
	}

	// Проверяем, удаляем ли вопрос из блока
	isBlockQuestion, _ := userState.Data["is_edit_block_questions"].(bool)

	if isBlockQuestion {
		// Удаление вопроса из блока
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		checklist := checklistData.(*types.BlockedCheckList)
		blockIdx, _ := userState.Data["current_block_index"].(int)

		// Удаляем вопрос из блока
		block := &checklist.Blocks[blockIdx]
		block.Questions = append(block.Questions[:idx], block.Questions[idx+1:]...)

		// Обновляем количество
		total := len(block.Questions)
		userState.Data["total_questions"] = total

		// Если вопросов не осталось - возвращаемся в редактор блока
		if total == 0 {
			delete(userState.Data, "is_edit_block_questions")
			delete(userState.Data, "edit_question_index")
			s.stateMgr.SetState(userID, userState)
			s.stateMgr.NavigateTo(userID, "block-editor")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			return
		}

		// Корректируем индекс если нужно
		if idx >= total {
			idx = total - 1
		}
		userState.Data["edit_question_index"] = idx

		// Возвращаемся к списку вопросов блока
		s.stateMgr.NavigateTo(userID, "block-view-questions")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

		log.Printf("[QuestionEditService] ✅ Вопрос %d удален из блока", idx+1)
	} else {
		// Удаление вопроса из простого чек-листа
		checklistData, _ := s.stateMgr.GetCheckList(userID)
		simpleChecklist := checklistData.(*types.SimpleCheckList)

		// Удаляем вопрос
		simpleChecklist.Questions = append(simpleChecklist.Questions[:idx], simpleChecklist.Questions[idx+1:]...)

		// Обновляем количество
		total := len(simpleChecklist.Questions)
		userState.Data["total_questions"] = total

		// Если вопросов не осталось - возвращаемся в редактор
		if total == 0 {
			delete(userState.Data, "edit_question_index")
			s.stateMgr.NavigateTo(userID, "simple-checklist-editor")
			s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
			return
		}

		// Корректируем индекс если нужно
		if idx >= total {
			idx = total - 1
		}
		userState.Data["edit_question_index"] = idx

		// Возвращаемся к списку вопросов
		s.stateMgr.NavigateTo(userID, "view-question")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

		log.Printf("[QuestionEditService] ✅ Вопрос %d удален", idx+1)
	}
}

// HandleCancelDelete отменяет удаление вопроса
func (s *QuestionEditService) HandleCancelDelete(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Возвращаемся к деталям вопроса
	s.stateMgr.NavigateTo(userID, "edit-question-detail")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}
