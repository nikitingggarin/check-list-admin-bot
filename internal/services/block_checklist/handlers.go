package block_checklist

import (
	"context"
	"fmt"
	"log"
	"strings"

	"telegram-bot/internal/formatters"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
	"telegram-bot/internal/storage/models"
	"telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ========== СОЗДАНИЕ ЧЕК-ЛИСТА ==========

func (s *BlockChecklistService) HandleCancelCreateChecklist(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendScreen(update.Message.Chat.ID, "admin-menu", userState)
}

func (s *BlockChecklistService) HandleCreateBlockChecklist(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	checklist := types.NewDraftBlockedCheckList(text)
	userState.SetBlockedCheckList(checklist)
	s.stateMgr.SetState(userID, userState)
	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendScreen(update.Message.Chat.ID, "block-checklist-editor", userState)
}

// ========== РАБОТА С БЛОКАМИ ==========

// HandleAddBlock начинает процесс создания нового блока
func (s *BlockChecklistService) HandleAddBlock(userID int64, update tgbotapi.Update, userState *state.UserState) {
	delete(userState.Data, "current_block_index")
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "edit-block-name")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] 🚀 Начало создания нового блока")
}

// HandleBlockNameInput обрабатывает ввод названия блока (и создание и редактирование)
func (s *BlockChecklistService) HandleBlockNameInput(userID int64, update tgbotapi.Update, userState *state.UserState, text string) {
	if len(strings.TrimSpace(text)) < 2 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Название блока должно содержать хотя бы 2 символа")
		return
	}

	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа")
		return
	}

	if blockIndex, exists := userState.Data["current_block_index"]; exists {
		idx, ok := blockIndex.(int)
		if !ok {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка: неверный индекс блока")
			return
		}

		if idx < 0 || idx >= len(checklist.Blocks) {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не найден")
			return
		}

		oldName := checklist.Blocks[idx].Name
		checklist.Blocks[idx].Name = text

		s.stateMgr.NavigateTo(userID, "block-editor")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

		log.Printf("[BlockChecklistService] ✅ Название блока изменено: '%s' → '%s' (индекс: %d)",
			oldName, text, idx)
	} else {
		block := types.NewBlock(text)
		checklist.AddBlock(block)

		blockIndex := len(checklist.Blocks) - 1
		userState.Data["current_block_index"] = blockIndex
		s.stateMgr.SetState(userID, userState)

		s.stateMgr.NavigateTo(userID, "block-editor")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

		log.Printf("[BlockChecklistService] ✅ Создан блок '%s' (индекс: %d)", text, blockIndex)
	}
}

// HandleBlockSelection обрабатывает выбор блока из списка
func (s *BlockChecklistService) HandleBlockSelection(userID int64, update tgbotapi.Update, userState *state.UserState, buttonText string) {
	blockIndex, ok := utils.ExtractBlockIndexFromButton(buttonText)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Не удалось определить выбранный блок")
		return
	}

	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа")
		return
	}

	if blockIndex < 0 || blockIndex >= len(checklist.Blocks) {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не найден")
		return
	}

	userState.Data["current_block_index"] = blockIndex
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "block-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	blockName := checklist.Blocks[blockIndex].Name
	log.Printf("[BlockChecklistService] ✅ Выбран блок '%s' (индекс: %d)", blockName, blockIndex)
}

// HandleCancelBlockEdit отменяет создание/редактирование блока
func (s *BlockChecklistService) HandleCancelBlockEdit(userID int64, update tgbotapi.Update, userState *state.UserState) {
	delete(userState.Data, "current_block_index")
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ❌ Отмена создания/редактирования блока")
}

// HandleBackFromBlockEditor возвращает из редактора блока к списку блоков
func (s *BlockChecklistService) HandleBackFromBlockEditor(userID int64, update tgbotapi.Update, userState *state.UserState) {
	delete(userState.Data, "current_block_index")
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ◀️ Возврат к списку блоков")
}

// HandleAddQuestionToBlock начинает добавление вопроса в текущий блок
func (s *BlockChecklistService) HandleAddQuestionToBlock(userID int64, update tgbotapi.Update, userState *state.UserState) {
	userState.Data["is_block_mode"] = true
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "select-question-type")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ➕ Начало добавления вопроса в блок")
}

// HandleEditBlockName начинает редактирование названия текущего блока
func (s *BlockChecklistService) HandleEditBlockName(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "edit-block-name")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ✏️ Начало редактирования названия блока")
}

// HandleBlockChecklistPreview показывает превью чек-листа с блоками
func (s *BlockChecklistService) HandleBlockChecklistPreview(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа")
		return
	}

	if len(checklist.Blocks) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист пустой. Добавьте блоки перед просмотром превью.")
		return
	}

	hasQuestions := false
	for _, block := range checklist.Blocks {
		if len(block.Questions) > 0 {
			hasQuestions = true
			break
		}
	}

	if !hasQuestions {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Все блоки пустые. Добавьте вопросы в блоки перед просмотром превью.")
		return
	}

	s.stateMgr.NavigateTo(userID, "checklist-preview")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] 👁️ Показано превью чек-листа с блоками: %s", checklist.Name)
}

// HandleEditChecklistTitle начинает редактирование названия чек-листа с блоками
func (s *BlockChecklistService) HandleEditChecklistTitle(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "edit-checklist-title")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleEditChecklistTitleInput обрабатывает ввод нового названия чек-листа с блоками
func (s *BlockChecklistService) HandleEditChecklistTitleInput(userID int64, update tgbotapi.Update, userState *state.UserState, newTitle string) {
	if len(strings.TrimSpace(newTitle)) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Название не может быть пустым")
		return
	}

	if len(newTitle) > 100 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Название слишком длинное (максимум 100 символов)")
		return
	}

	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа")
		return
	}

	oldName := checklist.Name
	checklist.Name = newTitle

	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ✅ Название чек-листа с блоками изменено: '%s' → '%s'", oldName, newTitle)
}

// HandleBackFromTitleEdit возвращает из редактирования названия
func (s *BlockChecklistService) HandleBackFromTitleEdit(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleBackFromPreview возвращает из превью
func (s *BlockChecklistService) HandleBackFromPreview(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleSaveDraft сохраняет черновик чек-листа с блоками
func (s *BlockChecklistService) HandleSaveDraft(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа для сохранения")
		return
	}

	if len(checklist.Blocks) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Нельзя сохранить чек-лист без блоков. Добавьте хотя бы один блок.")
		return
	}

	hasQuestions := false
	for _, block := range checklist.Blocks {
		if len(block.Questions) > 0 {
			hasQuestions = true
			break
		}
	}

	if !hasQuestions {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Нельзя сохранить чек-лист без вопросов. Добавьте вопросы в блоки.")
		return
	}

	ctx := context.Background()
	var savedChecklist *models.Checklist
	var err error

	if checklist.ID > 0 {
		log.Printf("[BlockChecklistService] Редактирование чек-листа с блоками ID=%d", checklist.ID)
		savedChecklist, err = s.checklistSvc.UpdateChecklist(ctx, checklist.ID, checklist, userID)
		if err != nil {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка обновления: "+err.Error())
			return
		}
		log.Printf("[BlockChecklistService] Старый чек-лист %d удален, создан новый ID=%d",
			checklist.ID, savedChecklist.ID)
	} else {
		log.Printf("[BlockChecklistService] Создание нового чек-листа с блоками")
		savedChecklist, err = s.checklistSvc.SaveBlockedChecklistDraft(ctx, checklist, userID)
		if err != nil {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка сохранения: "+err.Error())
			return
		}
	}

	totalQuestions := 0
	for _, block := range checklist.Blocks {
		totalQuestions += len(block.Questions)
	}

	var message string
	if checklist.ID > 0 {
		message = fmt.Sprintf("✅ Чек-лист с блоками обновлен!\n\n📋 Название: %s\n🔄 Старый ID: %d → Новый ID: %d\n🧱 Блоков: %d\n📊 Вопросов: %d\n\nЧек-лист доступен в разделе 'Мои чек-листы'",
			savedChecklist.Name, checklist.ID, savedChecklist.ID, len(checklist.Blocks), totalQuestions)
	} else {
		message = fmt.Sprintf("✅ Черновик сохранен!\n\n📋 Название: %s\n🏷️ ID: %d\n🧱 Блоков: %d\n📊 Вопросов: %d\n\nЧек-лист доступен в разделе 'Мои чек-листы'",
			savedChecklist.Name, savedChecklist.ID, len(checklist.Blocks), totalQuestions)
	}

	s.screenSvc.SendMessage(update.Message.Chat.ID, message)

	s.stateMgr.ClearCheckList(userID)
	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleSavePublish сохраняет и публикает чек-лист с блоками
func (s *BlockChecklistService) HandleSavePublish(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа для публикации")
		return
	}

	if len(checklist.Blocks) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Нельзя опубликовать чек-лист без блоков. Добавьте хотя бы один блок.")
		return
	}

	hasQuestions := false
	for _, block := range checklist.Blocks {
		if len(block.Questions) > 0 {
			hasQuestions = true
			break
		}
	}

	if !hasQuestions {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Нельзя опубликовать чек-лист без вопросов. Добавьте вопросы в блоки.")
		return
	}

	ctx := context.Background()
	var savedChecklist *models.Checklist
	var err error

	if checklist.ID > 0 {
		log.Printf("[BlockChecklistService] Редактирование и публикация чек-листа с блоками ID=%d", checklist.ID)
		savedChecklist, err = s.checklistSvc.UpdateChecklist(ctx, checklist.ID, checklist, userID)
		if err != nil {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка обновления: "+err.Error())
			return
		}
	} else {
		log.Printf("[BlockChecklistService] Создание и публикация нового чек-листа с блоками")
		savedChecklist, err = s.checklistSvc.SaveBlockedChecklistDraft(ctx, checklist, userID)
		if err != nil {
			s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка сохранения: "+err.Error())
			return
		}
	}

	// Подсчитываем общее количество вопросов ДО проверки ошибок
	totalQuestions := 0
	for _, block := range checklist.Blocks {
		totalQuestions += len(block.Questions)
	}

	// Публикуем чек-лист
	err = s.checklistSvc.PublishChecklist(ctx, savedChecklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "⚠️ Чек-лист сохранен как черновик, но не опубликован: "+err.Error())

		var msg string
		if checklist.ID > 0 {
			msg = fmt.Sprintf("🔄 Чек-лист с блоками обновлен как черновик\n\n📋 Название: %s\n🔄 Старый ID: %d → Новый ID: %d\n🧱 Блоков: %d\n📊 Вопросов: %d\n\nНо не опубликован из-за ошибки",
				savedChecklist.Name, checklist.ID, savedChecklist.ID, len(checklist.Blocks), totalQuestions)
		} else {
			msg = fmt.Sprintf("✅ Черновик сохранен!\n\n📋 Название: %s\n🏷️ ID: %d\n🧱 Блоков: %d\n📊 Вопросов: %d\n\nНо не опубликован из-за ошибки",
				savedChecklist.Name, savedChecklist.ID, len(checklist.Blocks), totalQuestions)
		}
		s.screenSvc.SendMessage(update.Message.Chat.ID, msg)
	} else {
		var message string
		if checklist.ID > 0 {
			message = fmt.Sprintf("🚀 Чек-лист с блоками обновлен и опубликован!\n\n📋 Название: %s\n🔄 Старый ID: %d → Новый ID: %d\n🧱 Блоков: %d\n📊 Вопросов: %d\n\nТеперь чек-лист доступен для прохождения",
				savedChecklist.Name, checklist.ID, savedChecklist.ID, len(checklist.Blocks), totalQuestions)
		} else {
			message = fmt.Sprintf("🚀 Чек-лист с блоками опубликован!\n\n📋 Название: %s\n🏷️ ID: %d\n🧱 Блоков: %d\n📊 Вопросов: %d\n\nТеперь чек-лист доступен для прохождения",
				savedChecklist.Name, savedChecklist.ID, len(checklist.Blocks), totalQuestions)
		}
		s.screenSvc.SendMessage(update.Message.Chat.ID, message)
	}

	// Очищаем состояние и возвращаемся в главное меню
	s.stateMgr.ClearCheckList(userID)
	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleConfirmExit начинает процесс подтверждения выхода в главное меню
func (s *BlockChecklistService) HandleConfirmExit(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "confirm-exit-block-checklist")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}

// HandleConfirmExitYes подтверждает выход в главное меню
func (s *BlockChecklistService) HandleConfirmExitYes(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.ClearCheckList(userID)
	s.stateMgr.NavigateTo(userID, "admin-menu")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	log.Printf("[BlockChecklistService] ✅ Выход из редактора блоков в главное меню")
}

// HandleConfirmExitNo отменяет выход в главное меню
func (s *BlockChecklistService) HandleConfirmExitNo(userID int64, update tgbotapi.Update, userState *state.UserState) {
	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	log.Printf("[BlockChecklistService] ❌ Отмена выхода из редактора блоков")
}

// HandleEditBlockQuestions начинает редактирование вопросов в блоке
func (s *BlockChecklistService) HandleEditBlockQuestions(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа")
		return
	}

	blockIdx, ok := userState.Data["current_block_index"].(int)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не выбран")
		return
	}

	if blockIdx < 0 || blockIdx >= len(checklist.Blocks) {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не найден")
		return
	}

	block := checklist.Blocks[blockIdx]

	if len(block.Questions) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ В блоке нет вопросов для редактирования")
		return
	}

	userState.Data["is_edit_block_questions"] = true
	userState.Data["edit_question_index"] = 0
	userState.Data["total_questions"] = len(block.Questions)
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "block-view-questions")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ✏️ Начало редактирования вопросов в блоке '%s'", block.Name)
}

// HandleBlockPreview показывает превью блока (без перехода на отдельный экран)
func (s *BlockChecklistService) HandleBlockPreview(userID int64, update tgbotapi.Update, userState *state.UserState) {
	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа")
		return
	}

	blockIdx, ok := userState.Data["current_block_index"].(int)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не выбран")
		return
	}

	if blockIdx < 0 || blockIdx >= len(checklist.Blocks) {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не найден")
		return
	}

	block := checklist.Blocks[blockIdx]

	if len(block.Questions) == 0 {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок пустой. Добавьте вопросы перед просмотром превью.")
		return
	}

	message := formatters.FormatBlockPreview(block)

	s.screenSvc.SendMessage(update.Message.Chat.ID, message)

	log.Printf("[BlockChecklistService] 👁️ Показано превью блока '%s'", block.Name)
}

// HandleBackFromBlockQuestions возвращает из списка вопросов блока
func (s *BlockChecklistService) HandleBackFromBlockQuestions(userID int64, update tgbotapi.Update, userState *state.UserState) {
	delete(userState.Data, "is_edit_block_questions")
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "block-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ◀️ Возврат из списка вопросов блока")
}

// HandleEditBlockQuestion начинает редактирование конкретного вопроса в блоке
func (s *BlockChecklistService) HandleEditBlockQuestion(userID int64, update tgbotapi.Update, userState *state.UserState, buttonText string) {
	questionNumber, err := utils.ExtractNumberFromPencilButton(buttonText)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Не удалось определить номер вопроса")
		return
	}

	checklistData, exists := s.stateMgr.GetCheckList(userID)
	if !exists {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Чек-лист не найден")
		return
	}

	checklist, ok := checklistData.(*types.BlockedCheckList)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Неверный тип чек-листа")
		return
	}

	blockIdx, ok := userState.Data["current_block_index"].(int)
	if !ok {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не выбран")
		return
	}

	if blockIdx < 0 || blockIdx >= len(checklist.Blocks) {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Блок не найден")
		return
	}

	block := checklist.Blocks[blockIdx]

	if questionNumber < 1 || questionNumber > len(block.Questions) {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Вопрос с таким номером не существует")
		return
	}

	userState.Data["edit_question_index"] = questionNumber - 1
	userState.Data["is_edit_block_questions"] = true
	s.stateMgr.SetState(userID, userState)

	s.stateMgr.NavigateTo(userID, "edit-question-detail")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[BlockChecklistService] ✅ Выбран вопрос %d в блоке '%s' для редактирования",
		questionNumber, block.Name)
}
