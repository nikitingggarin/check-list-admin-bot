package my_checklists

import (
	"fmt"
	"log"

	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/storage/models"
	"telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleChecklistNumber обрабатывает выбор чек-листа по номеру кнопки
func (s *MyChecklistsService) HandleChecklistNumber(userID int64, update tgbotapi.Update, userState *state.UserState, buttonText string) {
	// Извлекаем номер из кнопки (формат: "✏️ 1", "✏️ 2")
	number, err := utils.ExtractNumberFromPencilButton(buttonText)
	if err != nil || number < 1 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Не удалось определить выбранный чек-лист")
		return
	}

	// Получаем список чек-листов из состояния
	checklists, ok := userState.Data["my_checklists"].([]models.Checklist)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Список чек-листов не найден")
		return
	}

	// Проверяем диапазон
	if number > len(checklists) {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист с таким номером не существует")
		return
	}

	// Получаем выбранный чек-лист (0-based индекс)
	checklist := checklists[number-1]

	// Загружаем полные данные чек-листа
	dbChecklist, blocks, questions, answerOptions, err := s.checklistSvc.GetChecklistByID(checklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при загрузке чек-листа: "+err.Error())
		return
	}

	// ЗАГРУЖАЕМ ШАБЛОНЫ ДЛЯ ГРУППИРОВКИ ВОПРОСОВ ПО БЛОКАМ
	templates, err := s.checklistSvc.GetTemplatesByChecklistID(checklist.ID)
	if err != nil {
		log.Printf("[MyChecklistsService] ⚠️ Не удалось загрузить шаблоны для чек-листа %d: %v", checklist.ID, err)
		// Продолжаем без шаблонов
	}

	// Определяем тип чек-листа
	hasBlocks := len(blocks) > 0
	totalQuestions := len(questions)

	// Сохраняем данные в состояние
	userState.Data["current_checklist"] = dbChecklist
	userState.Data["has_blocks"] = hasBlocks
	userState.Data["total_questions"] = totalQuestions
	userState.Data["checklist_blocks"] = blocks
	userState.Data["checklist_questions"] = questions
	userState.Data["checklist_answer_options"] = answerOptions
	// Сохраняем шаблоны для группировки вопросов по блокам
	userState.Data["checklist_templates"] = templates

	s.stateMgr.SetState(userID, userState)

	// Переходим на экран деталей чек-листа (там будет превью + кнопки)
	s.stateMgr.NavigateTo(userID, "checklist-detail")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[MyChecklistsService] ✅ Загружен чек-лист ID=%d, тип=%v, вопросов=%d, шаблонов=%d",
		dbChecklist.ID, hasBlocks, totalQuestions, len(templates))
}

// HandleEditChecklist начинает редактирование чек-листа
func (s *MyChecklistsService) HandleEditChecklist(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklist, ok := userState.Data["current_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	hasBlocks, _ := userState.Data["has_blocks"].(bool)
	blocks, _ := userState.Data["checklist_blocks"].([]models.QuestionBlock)
	questions, _ := userState.Data["checklist_questions"].([]models.Question)
	answerOptions, _ := userState.Data["checklist_answer_options"].([]models.AnswerOption)

	if hasBlocks {
		// Редактирование чек-листа с блоками
		s.handleEditBlockChecklist(userID, update, userState, checklist, blocks, questions, answerOptions)
	} else {
		// Редактирование простого чек-листа
		s.handleEditSimpleChecklist(userID, update, userState, checklist, questions, answerOptions)
	}
}

// HandleDeleteChecklist начинает процесс удаления чек-листа
func (s *MyChecklistsService) HandleDeleteChecklist(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "confirm-delete-checklist")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleConfirmDelete подтверждает удаление чек-листа
func (s *MyChecklistsService) HandleConfirmDelete(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklist, ok := userState.Data["current_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	// Удаляем из базы данных
	err := s.checklistSvc.DeleteChecklist(checklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при удалении: "+err.Error())
		return
	}

	log.Printf("[MyChecklistsService] Чек-лист %d удален", checklist.ID)

	// Форматируем дату создания
	createdAtStr := "(дата недоступна)"
	if !checklist.CreatedAt.IsZero() && checklist.CreatedAt.Year() > 1 {
		createdAtStr = checklist.CreatedAt.Format("02.01.2006")
	}

	// Отправляем сообщение об успехе
	message := fmt.Sprintf("✅ Чек-лист удален:\n\n📋 Название: %s\n🏷️ ID: %d\n📅 Создан: %s\n\nОбновляю список чек-листов...",
		checklist.Name, checklist.ID, createdAtStr)
	s.screenSvc.SendMessage(update.Message.Chat.ID, message)

	// Очищаем данные удаленного чек-листа
	delete(userState.Data, "current_checklist")
	delete(userState.Data, "has_blocks")
	delete(userState.Data, "total_questions")
	delete(userState.Data, "checklist_blocks")
	delete(userState.Data, "checklist_questions")
	delete(userState.Data, "checklist_answer_options")
	s.stateMgr.SetState(userID, userState)

	// Обновляем список чек-листов
	s.refreshChecklistsList(userID, update, userState)
}

// HandleCancelDelete отменяет удаление чек-листа
func (s *MyChecklistsService) HandleCancelDelete(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "checklist-detail")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleBackFromDetail возвращает из деталей к списку
func (s *MyChecklistsService) HandleBackFromDetail(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Обновляем список чек-листов перед возвратом
	s.refreshChecklistsListOnBack(userID, update, userState)
}

// HandleBackFromList возвращает из списка в главное меню
func (s *MyChecklistsService) HandleBackFromList(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Очищаем список чек-листов
	delete(userState.Data, "my_checklists")
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandlePublishChecklist публикует чек-лист
func (s *MyChecklistsService) HandlePublishChecklist(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklist, ok := userState.Data["current_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	err := s.checklistSvc.PublishChecklist(checklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при публикации: "+err.Error())
		return
	}

	// Форматируем дату создания
	createdAtStr := "(дата недоступна)"
	if !checklist.CreatedAt.IsZero() && checklist.CreatedAt.Year() > 1 {
		createdAtStr = checklist.CreatedAt.Format("02.01.2006")
	}

	// Отправляем сообщение об успехе
	message := fmt.Sprintf("🚀 Чек-лист опубликован:\n\n📋 Название: %s\n🏷️ ID: %d\n📅 Создан: %s\n\nТеперь чек-лист доступен в разделе 'Опубликованные'",
		checklist.Name, checklist.ID, createdAtStr)
	s.screenSvc.SendMessage(update.Message.Chat.ID, message)

	// Обновляем список чек-листов
	s.refreshChecklistsList(userID, update, userState)
}
