package question_edit

import (
	"telegram-bot/internal/services/answers"
	"telegram-bot/internal/services/question"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
)

type QuestionEditService struct {
	stateMgr    manager.StateManager
	screenSvc   *screen.ScreenService
	questionSvc *question.QuestionService
	answersSvc  *answers.AnswersService
}

func NewQuestionEditService(stateMgr manager.StateManager, screenSvc *screen.ScreenService, questionSvc *question.QuestionService, answersSvc *answers.AnswersService) *QuestionEditService {
	return &QuestionEditService{
		stateMgr:    stateMgr,
		screenSvc:   screenSvc,
		questionSvc: questionSvc,
	}
}
