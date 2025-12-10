package answers

import (
	"telegram-bot/internal/services/question"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
)

type AnswersService struct {
	stateMgr    manager.StateManager
	screenSvc   *screen.ScreenService
	questionSvc *question.QuestionService
}

func NewAnswersService(stateMgr manager.StateManager, screenSvc *screen.ScreenService, questionSvc *question.QuestionService) *AnswersService {
	return &AnswersService{
		stateMgr:    stateMgr,
		screenSvc:   screenSvc,
		questionSvc: questionSvc,
	}
}
