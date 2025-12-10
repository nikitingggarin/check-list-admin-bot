package routes

import (
	"context"

	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
	"telegram-bot/internal/storage/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Router маршрутизирует сообщения к обработчикам
type Router struct {
	stateMgr             manager.StateManager
	userSvc              *service.UserService
	authRoute            *AuthRoutes
	menuRoute            *MenuRoutes
	simpleChecklistRoute *SimpleChecklistRoutes
	blockChecklistRoute  *BlockChecklistRoutes
	questionRoute        *QuestionRoutes
	answersRoute         *AnswersRoutes
	questionEditRoute    *QuestionEditRoutes
	myChecklistsRoute    *MyChecklistsRoutes
	publishedRoute       *PublishedChecklistsRoutes
}

// NewRouter создает новый роутер
func NewRouter(
	stateMgr manager.StateManager,
	userSvc *service.UserService,
	authRoute *AuthRoutes,
	menuRoute *MenuRoutes,
	simpleChecklistRoute *SimpleChecklistRoutes,
	blockChecklistRoute *BlockChecklistRoutes,
	questionRoute *QuestionRoutes,
	answersRoute *AnswersRoutes,
	questionEditRoute *QuestionEditRoutes,
	myChecklistsRoute *MyChecklistsRoutes,
	publishedRoute *PublishedChecklistsRoutes,
) *Router {
	return &Router{
		stateMgr:             stateMgr,
		userSvc:              userSvc,
		authRoute:            authRoute,
		menuRoute:            menuRoute,
		simpleChecklistRoute: simpleChecklistRoute,
		blockChecklistRoute:  blockChecklistRoute,
		questionRoute:        questionRoute,
		answersRoute:         answersRoute,
		questionEditRoute:    questionEditRoute,
		myChecklistsRoute:    myChecklistsRoute, // Инициализируем новый роут
		publishedRoute:       publishedRoute,    // ИНИЦИАЛИЗИРУЙ
	}
}

// Route определяет текущий экран и вызывает соответствующий обработчик
func (r *Router) Route(ctx context.Context, userID int64, update tgbotapi.Update, text string) {

	// Получаем состояние пользователя
	userState, exists := r.stateMgr.GetState(userID)
	if !exists {
		// Создаем новое состояние
		userState = state.NewUserState(nil, "authorize-admin")
		r.stateMgr.SetState(userID, userState)
	}

	// Получаем текущий экран
	currentScreen := userState.GetCurrentScreen()
	if currentScreen == "" {
		currentScreen = "authorize-admin"
		userState.SetCurrentScreen(currentScreen)
	}

	// Маршрутизация по экранам
	switch currentScreen {
	case "authorize-admin":
		r.authRoute.Route(ctx, userID, update, text, userState)
	case "admin-menu":
		r.menuRoute.Route(ctx, userID, update, text, userState)
	case "create-simple-checklist-name":
		r.simpleChecklistRoute.Route(ctx, userID, update, text, userState)
	case "create-block-checklist-name":
		r.blockChecklistRoute.Route(ctx, userID, update, text, userState)

	// ========== ЭКРАНЫ ДЛЯ ОПУБЛИКОВАННЫХ/ОТМЕНЕННЫХ ==========
	case "published-checklists-list":
		r.publishedRoute.Route(ctx, userID, update, text, userState)
	case "published-checklist-detail":
		r.publishedRoute.Route(ctx, userID, update, text, userState)

	// ========== ЭКРАНЫ ДЛЯ МОИХ ЧЕК-ЛИСТОВ ==========
	case "my-checklists-list":
		r.myChecklistsRoute.Route(ctx, userID, update, text, userState)
	case "checklist-detail":
		r.myChecklistsRoute.Route(ctx, userID, update, text, userState)
	case "confirm-delete-checklist":
		r.myChecklistsRoute.Route(ctx, userID, update, text, userState)

	// ========== ЭКРАНЫ ДЛЯ ПРОСТЫХ ЧЕК-ЛИСТОВ ==========
	case "simple-checklist-editor":
		r.simpleChecklistRoute.Route(ctx, userID, update, text, userState)
	case "edit-checklist-title":
		// Определяем тип чек-листа для выбора правильного роутера
		if userState.HasCheckList() {
			switch userState.CurrentCheckList.(type) {
			case *types.BlockedCheckList:
				r.blockChecklistRoute.Route(ctx, userID, update, text, userState)
			default:
				r.simpleChecklistRoute.Route(ctx, userID, update, text, userState)
			}
		} else {
			r.simpleChecklistRoute.Route(ctx, userID, update, text, userState)
		}
	case "checklist-preview":
		// Определяем тип чек-листа для выбора правильного роутера
		if userState.HasCheckList() {
			switch userState.CurrentCheckList.(type) {
			case *types.BlockedCheckList:
				r.blockChecklistRoute.Route(ctx, userID, update, text, userState)
			default:
				r.simpleChecklistRoute.Route(ctx, userID, update, text, userState)
			}
		} else {
			r.simpleChecklistRoute.Route(ctx, userID, update, text, userState)
		}
	case "confirm-exit-to-main-menu":
		r.simpleChecklistRoute.Route(ctx, userID, update, text, userState)

	// ========== ЭКРАНЫ ДЛЯ ЧЕК-ЛИСТОВ С БЛОКАМИ ==========
	case "block-checklist-editor":
		r.blockChecklistRoute.Route(ctx, userID, update, text, userState)
	case "edit-block-name":
		r.blockChecklistRoute.Route(ctx, userID, update, text, userState)
	case "block-editor":
		r.blockChecklistRoute.Route(ctx, userID, update, text, userState)
	case "confirm-exit-block-checklist":
		r.blockChecklistRoute.Route(ctx, userID, update, text, userState)
	case "block-view-questions":
		r.blockChecklistRoute.Route(ctx, userID, update, text, userState)

	// ========== ЭКРАНЫ ВОПРОСОВ ==========
	case "select-question-type":
		r.questionRoute.Route(ctx, userID, update, text, userState)
	case "enter-question-text":
		r.questionRoute.Route(ctx, userID, update, text, userState)

	// ========== ЭКРАНЫ ОТВЕТОВ ==========
	case "enter-answer-options":
		r.answersRoute.Route(ctx, userID, update, text, userState)
	case "select-correct-answers":
		r.answersRoute.Route(ctx, userID, update, text, userState)

	// ========== ЭКРАНЫ РЕДАКТИРОВАНИЯ ВОПРОСОВ ==========
	case "view-question":
		r.questionEditRoute.Route(ctx, userID, update, text, userState)
	case "edit-question-text":
		r.questionEditRoute.Route(ctx, userID, update, text, userState)
	case "edit-question-type":
		r.questionEditRoute.Route(ctx, userID, update, text, userState)
	case "confirm-delete-question":
		r.questionEditRoute.Route(ctx, userID, update, text, userState)
	case "edit-question-detail":
		r.questionEditRoute.Route(ctx, userID, update, text, userState)

	default:
		r.stateMgr.NavigateTo(userID, "authorize-admin")
	}
}
