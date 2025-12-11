package simple_checklist

import (
	"fmt"
	"log"
	"strings"

	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
	"telegram-bot/internal/storage/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ========== Методы для ввода названия ==========

// HandleCancelCreateSimpleChecklistName - отмена создания названия
func (r *SimpleChecklistService) HandleCancelCreateSimpleChecklistName(userID int64, update tgbotapi.Update, userState *state.UserState) {
	r.stateMgr.NavigateTo(userID, "admin-menu")
	r.screenSvc.SendScreen(update.Message.Chat.ID, "admin-menu", userState)
}

// HandleCreateSimpleChecklistName - создание названия простого чек-листа
func (r *SimpleChecklistService) HandleCreateSimpleChecklistName(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	checklist := types.NewDraftSimpleCheckList(text)
	userState.SetSimpleCheckList(checklist)
	r.stateMgr.SetState(userID, userState)
	r.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	r.screenSvc.SendScreen(update.Message.Chat.ID, "simple-checklist-editor", userState)
}

// ========== Методы для редактора чек-листа ==========

// HandleBtnBackToMainMenu - возврат в главное меню
func (r *SimpleChecklistService) HandleBtnBackToMainMenu(userID int64, update tgbotapi.Update, userState *state.UserState) {
	r.stateMgr.NavigateTo(userID, "admin-menu")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleBtnAddQuestion - добавление вопроса
func (r *SimpleChecklistService) HandleBtnAddQuestion(userID int64, update tgbotapi.Update, userState *state.UserState) {
	r.stateMgr.NavigateTo(userID, "select-question-type")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleEditTitle - начало редактирования названия
func (r *SimpleChecklistService) HandleEditTitle(userID int64, update tgbotapi.Update, userState *state.UserState) {
	r.stateMgr.NavigateTo(userID, "edit-checklist-title")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleNewTitleInput - обработка ввода нового названия
func (r *SimpleChecklistService) HandleNewTitleInput(userID int64, update tgbotapi.Update, userState *state.UserState, newTitle string) {
	// Проверяем, что название не пустое
	if len(strings.TrimSpace(newTitle)) == 0 {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Название не может быть пустым")
		return
	}

	// Проверяем максимальную длину
	if len(newTitle) > 100 {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Название слишком длинное (максимум 100 символов)")
		return
	}

	// Получаем текущий чек-лист
	checklistData, exists := r.stateMgr.GetCheckList(userID)
	if !exists {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	// Меняем название в зависимости от типа чек-листа
	switch checklist := checklistData.(type) {
	case *types.SimpleCheckList:
		oldName := checklist.Name
		checklist.Name = newTitle
		log.Printf("[SimpleChecklistService] ✅ Название простого чек-листа изменено: '%s' → '%s'", oldName, newTitle)
	}

	// Возвращаемся в редактор
	r.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleEditQuestions - начало редактирования вопросов
func (r *SimpleChecklistService) HandleEditQuestions(userID int64, update tgbotapi.Update, userState *state.UserState) {
	r.questionEditSvc.HandleEditQuestions(userID, update, userState)
}

// HandlePreview - показ превью чек-листа
func (r *SimpleChecklistService) HandlePreview(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Проверяем, есть ли вопросы в чек-листе
	checklistData, exists := r.stateMgr.GetCheckList(userID)
	if !exists {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	simpleChecklist, ok := checklistData.(*types.SimpleCheckList)
	if !ok {
		// Если это не SimpleCheckList - показываем сообщение
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Превью для этого типа чек-листа недоступно")
		return
	}

	hasQuestions := len(simpleChecklist.Questions) > 0

	if !hasQuestions {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист пустой. Добавьте вопросы перед просмотром превью.")
		return
	}

	// Переходим к экрану превью
	r.stateMgr.NavigateTo(userID, "checklist-preview")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleBack - обработка кнопки "Назад"
func (r *SimpleChecklistService) HandleBack(userID int64, update tgbotapi.Update, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	switch currentScreen {
	case "edit-checklist-title":
		// Возвращаемся из редактирования названия в редактор
		r.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	case "checklist-preview":
		// Возвращаемся из превью в редактор
		r.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	case "confirm-exit-to-main-menu":
		// Возвращаемся из подтверждения выхода в редактор
		r.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	default:
		r.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	}

	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleSaveDraft - сохранение черновика
func (r *SimpleChecklistService) HandleSaveDraft(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Получаем текущий чек-лист
	checklistData, exists := r.stateMgr.GetCheckList(userID)
	if !exists {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	simpleChecklist, ok := checklistData.(*types.SimpleCheckList)
	if !ok {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа для сохранения")
		return
	}

	// Проверяем, есть ли вопросы
	if len(simpleChecklist.Questions) == 0 {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Нельзя сохранить пустой чек-лист. Добавьте хотя бы один вопрос.")
		return
	}

	var savedChecklist *models.Checklist
	var err error

	// Проверяем, редактируем ли существующий чек-лист (ID > 0)
	if simpleChecklist.ID > 0 {
		// РЕДАКТИРОВАНИЕ: удаляем старый и создаем новый
		log.Printf("[SimpleChecklistService] Редактирование чек-листа ID=%d", simpleChecklist.ID)
		savedChecklist, err = r.checklistSvc.UpdateChecklist(simpleChecklist.ID, simpleChecklist, userID)
		if err != nil {
			r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка обновления: "+err.Error())
			return
		}
		log.Printf("[SimpleChecklistService] Старый чек-лист %d удален, создан новый ID=%d",
			simpleChecklist.ID, savedChecklist.ID)
	} else {
		// СОЗДАНИЕ: просто создаем новый
		log.Printf("[SimpleChecklistService] Создание нового чек-листа")
		savedChecklist, err = r.checklistSvc.SaveSimpleChecklistDraft(simpleChecklist, userID)
		if err != nil {
			r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка сохранения: "+err.Error())
			return
		}
	}

	// Отправляем сообщение об успехе с указанием нового ID
	var message string
	if simpleChecklist.ID > 0 {
		message = fmt.Sprintf("✅ Чек-лист обновлен!\n\n📋 Название: %s\n🔄 Старый ID: %d → Новый ID: %d\n📊 Вопросов: %d\n\nЧек-лист доступен в разделе 'Мои чек-листы'",
			savedChecklist.Name, simpleChecklist.ID, savedChecklist.ID, len(simpleChecklist.Questions))
	} else {
		message = fmt.Sprintf("✅ Черновик сохранен!\n\n📋 Название: %s\n🏷️ ID: %d\n📊 Вопросов: %d\n\nЧек-лист доступен в разделе 'Мои чек-листы'",
			savedChecklist.Name, savedChecklist.ID, len(simpleChecklist.Questions))
	}

	r.screenSvc.SendMessage(update.Message.Chat.ID, message)

	// Очищаем состояние и возвращаемся в главное меню
	r.stateMgr.ClearCheckList(userID)
	r.stateMgr.NavigateTo(userID, "admin-menu")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleSavePublish - сохранение и публикация
func (r *SimpleChecklistService) HandleSavePublish(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Получаем текущий чек-лист
	checklistData, exists := r.stateMgr.GetCheckList(userID)
	if !exists {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	simpleChecklist, ok := checklistData.(*types.SimpleCheckList)
	if !ok {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа для публикации")
		return
	}

	// Проверяем, есть ли вопросы
	if len(simpleChecklist.Questions) == 0 {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Нельзя опубликовать пустой чек-лист. Добавьте хотя бы один вопрос.")
		return
	}

	var savedChecklist *models.Checklist
	var err error

	// Проверяем, редактируем ли существующий чек-лист (ID > 0)
	if simpleChecklist.ID > 0 {
		// РЕДАКТИРОВАНИЕ: удаляем старый и создаем новый
		log.Printf("[SimpleChecklistService] Редактирование и публикация чек-листа ID=%d", simpleChecklist.ID)
		savedChecklist, err = r.checklistSvc.UpdateChecklist(simpleChecklist.ID, simpleChecklist, userID)
		if err != nil {
			r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка обновления: "+err.Error())
			return
		}
	} else {
		// СОЗДАНИЕ: просто создаем новый
		log.Printf("[SimpleChecklistService] Создание и публикация нового чек-листа")
		savedChecklist, err = r.checklistSvc.SaveSimpleChecklistDraft(simpleChecklist, userID)
		if err != nil {
			r.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка сохранения: "+err.Error())
			return
		}
	}

	// Публикуем чек-лист
	err = r.checklistSvc.PublishChecklist(savedChecklist.ID)
	if err != nil {
		r.screenSvc.SendMessage(update.Message.Chat.ID, "⚠️ Чек-лист сохранен как черновик, но не опубликован: "+err.Error())

		// Показываем сообщение с новым ID
		var msg string
		if simpleChecklist.ID > 0 {
			msg = fmt.Sprintf("🔄 Чек-лист обновлен как черновик\n\n📋 Название: %s\n🔄 Старый ID: %d → Новый ID: %d\n📊 Вопросов: %d\n\nНо не опубликован из-за ошибки",
				savedChecklist.Name, simpleChecklist.ID, savedChecklist.ID, len(simpleChecklist.Questions))
		} else {
			msg = fmt.Sprintf("✅ Черновик сохранен!\n\n📋 Название: %s\n🏷️ ID: %d\n📊 Вопросов: %d\n\nНо не опубликован из-за ошибки",
				savedChecklist.Name, savedChecklist.ID, len(simpleChecklist.Questions))
		}
		r.screenSvc.SendMessage(update.Message.Chat.ID, msg)
	} else {
		// Успешная публикация
		var message string
		if simpleChecklist.ID > 0 {
			message = fmt.Sprintf("🚀 Чек-лист обновлен и опубликован!\n\n📋 Название: %s\n🔄 Старый ID: %d → Новый ID: %d\n📊 Вопросов: %d\n\nТеперь чек-лист доступен для прохождения",
				savedChecklist.Name, simpleChecklist.ID, savedChecklist.ID, len(simpleChecklist.Questions))
		} else {
			message = fmt.Sprintf("🚀 Чек-лист опубликован!\n\n📋 Название: %s\n🏷️ ID: %d\n📊 Вопросов: %d\n\nТеперь чек-лист доступен для прохождения",
				savedChecklist.Name, savedChecklist.ID, len(simpleChecklist.Questions))
		}
		r.screenSvc.SendMessage(update.Message.Chat.ID, message)
	}

	// Очищаем состояние и возвращаемся в главное меню
	r.stateMgr.ClearCheckList(userID)
	r.stateMgr.NavigateTo(userID, "admin-menu")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleConfirmExit - начало процесса подтверждения выхода
func (r *SimpleChecklistService) HandleConfirmExit(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Переходим к экрану подтверждения выхода
	r.stateMgr.NavigateTo(userID, "confirm-exit-to-main-menu")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleConfirmExitYes - подтверждение выхода в главное меню
func (r *SimpleChecklistService) HandleConfirmExitYes(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Очищаем чек-лист
	r.stateMgr.ClearCheckList(userID)

	// Переходим в главное меню
	r.stateMgr.NavigateTo(userID, "admin-menu")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[SimpleChecklistService] ✅ Пользователь %d вышел в главное меню", userID)
}

// HandleConfirmExitNo - отмена выхода в главное меню
func (r *SimpleChecklistService) HandleConfirmExitNo(userID int64, update tgbotapi.Update, userState *state.UserState) {
	// Возвращаемся в редактор чек-листа
	r.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[SimpleChecklistService] ❌ Пользователь %d отменил выход", userID)
}
