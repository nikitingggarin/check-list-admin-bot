package simple_checklist

import (
	"telegram-bot/internal/services/question_edit"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/storage/service"
)

// SimpleChecklistService - объединенный сервис для работы с простыми чек-листами
type SimpleChecklistService struct {
	stateMgr        manager.StateManager
	screenSvc       *screen.ScreenService
	questionEditSvc *question_edit.QuestionEditService
	checklistSvc    *service.ChecklistService
}

func NewSimpleChecklistService(
	stateMgr manager.StateManager,
	screenSvc *screen.ScreenService,
	questionEditSvc *question_edit.QuestionEditService,
	checklistSvc *service.ChecklistService,
) *SimpleChecklistService {
	return &SimpleChecklistService{
		stateMgr:        stateMgr,
		screenSvc:       screenSvc,
		questionEditSvc: questionEditSvc,
		checklistSvc:    checklistSvc,
	}
}
