package menu

import (
	"telegram-bot/internal/services/published_checklists"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/storage/service"
)

// MenuService содержит методы для маршрутизации главного меню
type MenuService struct {
	stateMgr     manager.StateManager
	screenSvc    *screen.ScreenService
	checklistSvc *service.ChecklistService // Добавляем ChecklistService
	publishedSvc *published_checklists.PublishedChecklistsService
}

// NewMenuService создает новый роутер для главного меню
func NewMenuService(
	stateMgr manager.StateManager,
	screenSvc *screen.ScreenService,
	checklistSvc *service.ChecklistService, // Добавляем параметр
	publishedSvc *published_checklists.PublishedChecklistsService,
) *MenuService {
	return &MenuService{
		stateMgr:     stateMgr,
		screenSvc:    screenSvc,
		checklistSvc: checklistSvc,
		publishedSvc: publishedSvc, // ИНИЦИАЛИЗИРУЙ
	}
}
