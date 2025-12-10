package routes

import (
	"context"
	"log"
	"strings"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/block_checklist"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BlockChecklistRoutes содержит методы для маршрутизации чек-листов с блоками
type BlockChecklistRoutes struct {
	checklistSvc *block_checklist.BlockChecklistService
	screenSvc    *screen.ScreenService
}

func NewBlockChecklistRoutes(checklistSvc *block_checklist.BlockChecklistService, screenSvc *screen.ScreenService) *BlockChecklistRoutes {
	return &BlockChecklistRoutes{
		checklistSvc: checklistSvc,
		screenSvc:    screenSvc,
	}
}

// Route маршрутизирует сообщение в контексте чек-листов с блоками
func (r *BlockChecklistRoutes) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	log.Printf("[BlockChecklistRoutes] 👤 UserID: %d | 💬 Текст: %s | Экран: %s", userID, text, currentScreen)

	switch currentScreen {
	case "create-block-checklist-name":
		r.handleCreateNameScreen(userID, update, text, userState)
	case "block-checklist-editor":
		r.handleBlockListScreen(userID, update, text, userState)
	case "edit-block-name":
		r.handleEditBlockNameScreen(userID, update, text, userState)
	case "block-editor":
		r.handleBlockEditorScreen(userID, update, text, userState)
	case "edit-checklist-title":
		r.handleEditChecklistTitleScreen(userID, update, text, userState)
	case "checklist-preview":
		r.handleChecklistPreviewScreen(userID, update, text, userState)
	case "confirm-exit-block-checklist":
		r.handleConfirmExitScreen(userID, update, text, userState)
	case "block-view-questions":
		r.handleBlockViewQuestionsScreen(userID, update, text, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// handleCreateNameScreen обрабатывает экран создания названия чек-листа
func (r *BlockChecklistRoutes) handleCreateNameScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnBack:
		r.checklistSvc.HandleCancelCreateChecklist(userID, update, userState)
	default:
		r.checklistSvc.HandleCreateBlockChecklist(userID, update, userState, text)
	}
}

// handleBlockListScreen обрабатывает главный экран редактора блоков
func (r *BlockChecklistRoutes) handleBlockListScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnAddBlock:
		r.checklistSvc.HandleAddBlock(userID, update, userState)
	case buttons.BtnPreview:
		r.checklistSvc.HandleBlockChecklistPreview(userID, update, userState)
	case buttons.BtnEditTitleChecklist:
		r.checklistSvc.HandleEditChecklistTitle(userID, update, userState)
	case buttons.BtnBackToMainMenu:
		r.checklistSvc.HandleConfirmExit(userID, update, userState)
	default:
		if isBlockButton(text) {
			r.checklistSvc.HandleBlockSelection(userID, update, userState, text)
		} else {
			r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
		}
	}
}

// handleEditBlockNameScreen обрабатывает экран ввода названия блока
func (r *BlockChecklistRoutes) handleEditBlockNameScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnBack:
		r.checklistSvc.HandleCancelBlockEdit(userID, update, userState)
	default:
		// Вызываем единый метод для обработки ввода названия блока
		r.checklistSvc.HandleBlockNameInput(userID, update, userState, text)
	}
}

// handleBlockEditorScreen обрабатывает экран редактора конкретного блока
func (r *BlockChecklistRoutes) handleBlockEditorScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnBackToBlockList:
		r.checklistSvc.HandleBackFromBlockEditor(userID, update, userState)
	case buttons.BtnAddQuestion:
		r.checklistSvc.HandleAddQuestionToBlock(userID, update, userState)
	case buttons.BtnEditTitleBlockChecklist:
		r.checklistSvc.HandleEditBlockName(userID, update, userState)
	case buttons.BtnEditQuestionChecklist:
		r.checklistSvc.HandleEditBlockQuestions(userID, update, userState)
	case buttons.BtnPreviewBlock:
		r.checklistSvc.HandleBlockPreview(userID, update, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// handleEditChecklistTitleScreen обрабатывает экран редактирования названия чек-листа
func (r *BlockChecklistRoutes) handleEditChecklistTitleScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnBack:
		r.checklistSvc.HandleBackFromTitleEdit(userID, update, userState)
	default:
		r.checklistSvc.HandleEditChecklistTitleInput(userID, update, userState, text)
	}
}

// handleChecklistPreviewScreen обрабатывает экран превью чек-листа
func (r *BlockChecklistRoutes) handleChecklistPreviewScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnBack:
		r.checklistSvc.HandleBackFromPreview(userID, update, userState)
	case buttons.BtnSaveDraft:
		r.checklistSvc.HandleSaveDraft(userID, update, userState)
	case buttons.BtnSavePublish:
		r.checklistSvc.HandleSavePublish(userID, update, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// handleConfirmExitScreen обрабатывает экран подтверждения выхода
func (r *BlockChecklistRoutes) handleConfirmExitScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnYes:
		r.checklistSvc.HandleConfirmExitYes(userID, update, userState)
	case buttons.BtnNo:
		r.checklistSvc.HandleConfirmExitNo(userID, update, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// handleBlockViewQuestionsScreen обрабатывает экран списка вопросов блока
func (r *BlockChecklistRoutes) handleBlockViewQuestionsScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	// Проверяем, является ли текст кнопкой с карандашом и номером
	if utils.IsPencilNumberButton(text) {
		// Начинаем редактирование вопроса в блоке
		r.checklistSvc.HandleEditBlockQuestion(userID, update, userState, text)
		return
	}

	switch text {
	case buttons.BtnBack:
		r.checklistSvc.HandleBackFromBlockQuestions(userID, update, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// isBlockButton проверяет, является ли текст кнопкой блока
func isBlockButton(text string) bool {
	return strings.HasPrefix(text, "🧱") || strings.HasPrefix(text, "📭") || strings.HasPrefix(text, "🏗️")
}
