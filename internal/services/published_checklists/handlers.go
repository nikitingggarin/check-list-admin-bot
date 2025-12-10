package published_checklists

import (
	"context"
	"fmt"
	"log"

	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/storage/models"
	"telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandlePublishedChecklists показывает опубликованные чек-листы
func (s *PublishedChecklistsService) HandlePublishedChecklists(userID int64, update tgbotapi.Update, userState *state.UserState) {
	ctx := context.Background()
	checklists, err := s.checklistSvc.GetUserPublished(ctx, userID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при получении опубликованных чек-листов: "+err.Error())
		return
	}

	if len(checklists) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "🚀 У вас пока нет опубликованных чек-листов.")
		return
	}

	// Сохраняем список и тип (published)
	userState.Data["published_checklists"] = checklists
	userState.Data["checklists_type"] = "published" // тип списка: published/unpublished
	s.stateMgr.SetState(userID, userState)

	// Переходим на экран списка
	s.stateMgr.NavigateTo(userID, "published-checklists-list")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[PublishedChecklistsService] ✅ Пользователь %d просмотрел опубликованные чек-листы (%d шт.)", userID, len(checklists))
}

// HandleUnpublishedChecklists показывает отмененные чек-листы
func (s *PublishedChecklistsService) HandleUnpublishedChecklists(userID int64, update tgbotapi.Update, userState *state.UserState) {
	ctx := context.Background()
	checklists, err := s.checklistSvc.GetUserUnpublished(ctx, userID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при получении отмененных чек-листов: "+err.Error())
		return
	}

	if len(checklists) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "🚫 У вас пока нет отмененных публикаций.")
		return
	}

	// Сохраняем список и тип (unpublished)
	userState.Data["published_checklists"] = checklists
	userState.Data["checklists_type"] = "unpublished"
	s.stateMgr.SetState(userID, userState)

	// Переходим на экран списка
	s.stateMgr.NavigateTo(userID, "published-checklists-list")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[PublishedChecklistsService] ✅ Пользователь %d просмотрел отмененные чек-листы (%d шт.)", userID, len(checklists))
}

// HandleChecklistNumber обрабатывает выбор чек-листа по номеру
func (s *PublishedChecklistsService) HandleChecklistNumber(userID int64, update tgbotapi.Update, userState *state.UserState, buttonText string) {
	// Извлекаем номер из кнопки (формат: "✏️ 1", "✏️ 2")
	number, err := utils.ExtractNumberFromPencilButton(buttonText)
	if err != nil || number < 1 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Не удалось определить выбранный чек-лист")
		return
	}

	// Получаем список чек-листов из состояния
	checklists, ok := userState.Data["published_checklists"].([]models.Checklist)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Список чек-листов не найден")
		return
	}

	// Проверяем диапазон
	if number > len(checklists) {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист с таким номером не существует")
		return
	}

	// Получаем выбранный чек-лист
	checklist := checklists[number-1]
	checklistType, _ := userState.Data["checklists_type"].(string)

	// Загружаем полные данные чек-листа
	ctx := context.Background()
	dbChecklist, blocks, questions, answerOptions, err := s.checklistSvc.GetChecklistByID(ctx, checklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при загрузке чек-листа: "+err.Error())
		return
	}

	// ДОБАВЛЕНО: Загружаем шаблоны для группировки вопросов по блокам
	templates, err := s.checklistSvc.GetTemplatesByChecklistID(ctx, checklist.ID)
	if err != nil {
		log.Printf("[PublishedChecklistsService] ⚠️ Не удалось загрузить шаблоны для чек-листа %d: %v", checklist.ID, err)
		// Продолжаем без шаблонов
	}

	// Определяем тип чек-листа
	hasBlocks := len(blocks) > 0
	totalQuestions := len(questions)

	// Сохраняем данные в состояние
	userState.Data["current_published_checklist"] = dbChecklist
	userState.Data["published_has_blocks"] = hasBlocks
	userState.Data["published_total_questions"] = totalQuestions
	userState.Data["published_checklist_blocks"] = blocks
	userState.Data["published_checklist_questions"] = questions
	userState.Data["published_checklist_answer_options"] = answerOptions
	// ДОБАВЛЕНО: Сохраняем шаблоны
	userState.Data["published_checklist_templates"] = templates
	userState.Data["current_checklist_type"] = checklistType

	s.stateMgr.SetState(userID, userState)

	// Переходим на экран деталей
	s.stateMgr.NavigateTo(userID, "published-checklist-detail")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[PublishedChecklistsService] ✅ Загружен чек-лист ID=%d, тип=%s, вопросов=%d, шаблонов=%d",
		dbChecklist.ID, checklistType, totalQuestions, len(templates))
}

// HandleUnpublish снимает чек-лист с публикации
func (s *PublishedChecklistsService) HandleUnpublish(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklist, ok := userState.Data["current_published_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	ctx := context.Background()
	err := s.checklistSvc.UnpublishChecklist(ctx, checklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при снятии с публикации: "+err.Error())
		return
	}

	// Отправляем сообщение об успехе
	message := fmt.Sprintf("🚫 Чек-лист снят с публикации:\n\n📋 Название: %s\n🏷️ ID: %d\n\nТеперь он доступен в разделе 'Отмененные публикации'",
		checklist.Name, checklist.ID)
	s.screenSvc.SendMessage(update.Message.Chat.ID, message)

	// Возвращаемся в главное меню
	s.cleanupPublishedChecklistData(userState)
	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[PublishedChecklistsService] ✅ Чек-лист %d снят с публикации", checklist.ID)
}

// HandleRepublish возвращает чек-лист в публикацию
func (s *PublishedChecklistsService) HandleRepublish(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklist, ok := userState.Data["current_published_checklist"].(*models.Checklist)
	if !ok || checklist == nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	ctx := context.Background()
	err := s.checklistSvc.RepublishChecklist(ctx, checklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при возврате в публикацию: "+err.Error())
		return
	}

	// Отправляем сообщение об успехе
	message := fmt.Sprintf("🚀 Чек-лист возвращен в публикацию:\n\n📋 Название: %s\n🏷️ ID: %d\n\nТеперь он доступен в разделе 'Опубликованные'",
		checklist.Name, checklist.ID)
	s.screenSvc.SendMessage(update.Message.Chat.ID, message)

	// Возвращаемся в главное меню
	s.cleanupPublishedChecklistData(userState)
	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[PublishedChecklistsService] ✅ Чек-лист %d возвращен в публикацию", checklist.ID)
}

// HandleBackFromDetail возвращает из деталей к списку
func (s *PublishedChecklistsService) HandleBackFromDetail(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.cleanupPublishedChecklistData(userState)
	s.stateMgr.NavigateTo(userID, "published-checklists-list")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleBackFromList возвращает из списка в главное меню
func (s *PublishedChecklistsService) HandleBackFromList(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Очищаем данные
	delete(userState.Data, "published_checklists")
	delete(userState.Data, "checklists_type")
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}
