package auth

import (
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/storage/service"
)

type AuthService struct {
	stateMgr  manager.StateManager
	userSvc   *service.UserService
	screenSvc *screen.ScreenService
}

func NewAuthService(stateMgr manager.StateManager, userSvc *service.UserService, screenSvc *screen.ScreenService) *AuthService {
	return &AuthService{
		stateMgr:  stateMgr,
		userSvc:   userSvc,
		screenSvc: screenSvc,
	}
}
