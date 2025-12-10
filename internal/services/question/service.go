package question

import (
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
)

type QuestionService struct {
	stateMgr  manager.StateManager
	screenSvc *screen.ScreenService
}

func NewQuestionService(stateMgr manager.StateManager, screenSvc *screen.ScreenService) *QuestionService {
	return &QuestionService{
		stateMgr:  stateMgr,
		screenSvc: screenSvc,
	}
}
