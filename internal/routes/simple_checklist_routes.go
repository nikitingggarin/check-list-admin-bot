package routes

import (
	"context"
	"log"

	"telegram-bot/internal/buttons"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/services/simple_checklist"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SimpleChecklistRoutes содержит методы для маршрутизации простых чек-листов
type SimpleChecklistRoutes struct {
	simpleChecklistSvc *simple_checklist.SimpleChecklistService
	screenSvc          *screen.ScreenService
}

func NewSimpleChecklistRoutes(simpleChecklistSvc *simple_checklist.SimpleChecklistService, screenSvc *screen.ScreenService) *SimpleChecklistRoutes {
	return &SimpleChecklistRoutes{
		simpleChecklistSvc: simpleChecklistSvc,
		screenSvc:          screenSvc,
	}
}

// Route маршрутизирует сообщение в контексте простых чек-листов
func (r *SimpleChecklistRoutes) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	log.Printf("[SimpleChecklistRoutes] 👤 UserID: %d | 💬 Текст: %s | Экран: %s", userID, text, currentScreen)

	// Маршрутизация по экранам
	switch currentScreen {
	case "create-simple-checklist-name":
		r.handleCreateNameScreen(userID, update, text, userState)
	case "simple-checklist-editor",
		"block-checklist-editor",
		"edit-checklist-title",
		"checklist-preview",
		"confirm-exit-to-main-menu":
		r.handleEditorScreen(userID, update, text, userState)
	default:
		r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
	}
}

// handleCreateNameScreen обрабатывает экран создания названия
func (r *SimpleChecklistRoutes) handleCreateNameScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	switch text {
	case buttons.BtnBack:
		r.simpleChecklistSvc.HandleCancelCreateSimpleChecklistName(userID, update, userState)
	default:
		r.simpleChecklistSvc.HandleCreateSimpleChecklistName(userID, update, userState, text)
	}
}

// handleEditorScreen обрабатывает экраны редактора
func (r *SimpleChecklistRoutes) handleEditorScreen(userID int64, update tgbotapi.Update, text string, userState *state.UserState) {
	currentScreen := userState.GetCurrentScreen()

	switch text {
	case buttons.BtnBackToMainMenu:
		// Переходим к подтверждению выхода
		r.simpleChecklistSvc.HandleConfirmExit(userID, update, userState)
	case buttons.BtnAddQuestion:
		r.simpleChecklistSvc.HandleBtnAddQuestion(userID, update, userState)
	case buttons.BtnEditTitleChecklist:
		// Начало редактирования названия
		r.simpleChecklistSvc.HandleEditTitle(userID, update, userState)
	case buttons.BtnEditQuestionChecklist:
		// Редактирование вопросов
		r.simpleChecklistSvc.HandleEditQuestions(userID, update, userState)
	case buttons.BtnPreview:
		// Превью чек-листа
		r.simpleChecklistSvc.HandlePreview(userID, update, userState)
	case buttons.BtnBack:
		// Обработка кнопки "Назад"
		r.simpleChecklistSvc.HandleBack(userID, update, userState)
	case buttons.BtnSaveDraft:
		// Сохранение черновика из превью
		r.simpleChecklistSvc.HandleSaveDraft(userID, update, userState)
	case buttons.BtnSavePublish:
		// Сохранение и публикация из превью
		r.simpleChecklistSvc.HandleSavePublish(userID, update, userState)
	case buttons.BtnYes:
		// Подтверждение выхода в главное меню
		r.simpleChecklistSvc.HandleConfirmExitYes(userID, update, userState)
	case buttons.BtnNo:
		// Отмена выхода в главное меню
		r.simpleChecklistSvc.HandleConfirmExitNo(userID, update, userState)
	default:
		// Если находимся в режиме редактирования названия - обрабатываем ввод нового названия
		if currentScreen == "edit-checklist-title" {
			r.simpleChecklistSvc.HandleNewTitleInput(userID, update, userState, text)
		} else {
			r.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
		}
	}
}
