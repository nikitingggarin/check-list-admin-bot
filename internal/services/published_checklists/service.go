package published_checklists

import (
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/storage/service"
)

type PublishedChecklistsService struct {
	stateMgr     manager.StateManager
	screenSvc    *screen.ScreenService
	checklistSvc *service.ChecklistService
}

func NewPublishedChecklistsService(
	stateMgr manager.StateManager,
	screenSvc *screen.ScreenService,
	checklistSvc *service.ChecklistService,
) *PublishedChecklistsService {
	return &PublishedChecklistsService{
		stateMgr:     stateMgr,
		screenSvc:    screenSvc,
		checklistSvc: checklistSvc,
	}
}
