package my_checklists

import (
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/storage/service"
)

type MyChecklistsService struct {
	stateMgr     manager.StateManager
	screenSvc    *screen.ScreenService
	checklistSvc *service.ChecklistService
}

func NewMyChecklistsService(
	stateMgr manager.StateManager,
	screenSvc *screen.ScreenService,
	checklistSvc *service.ChecklistService,
) *MyChecklistsService {
	return &MyChecklistsService{
		stateMgr:     stateMgr,
		screenSvc:    screenSvc,
		checklistSvc: checklistSvc,
	}
}
