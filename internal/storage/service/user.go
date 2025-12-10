package service

import (
	"context"
	"fmt"
	"telegram-bot/internal/storage/models"
	"telegram-bot/internal/storage/repositories"
)

type UserService struct {
	user repositories.Repositories
}

func NewUserService(user repositories.Repositories) *UserService {
	return &UserService{user: user}
}

func (s *UserService) AuthorizeAdmin(ctx context.Context, telegramID int64) (*models.User, error) {
	user, err := s.user.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("у вас нет доступа к боту")
	}

	if user.Role != models.RoleAdmin {
		return nil, fmt.Errorf("доступ запрещен: требуется роль администратора")
	}

	return user, nil
}
