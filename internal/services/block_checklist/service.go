package block_checklist

import (
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/storage/service"
)

type BlockChecklistService struct {
	stateMgr     manager.StateManager
	screenSvc    *screen.ScreenService
	checklistSvc *service.ChecklistService
}

func NewBlockChecklistService(
	stateMgr manager.StateManager,
	screenSvc *screen.ScreenService,
	checklistSvc *service.ChecklistService,
) *BlockChecklistService {
	return &BlockChecklistService{
		stateMgr:     stateMgr,
		screenSvc:    screenSvc,
		checklistSvc: checklistSvc,
	}
}
