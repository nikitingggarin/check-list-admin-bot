package auth

import (
	"telegram-bot/internal/state_manager/types"
	"telegram-bot/internal/storage/models"
)

func convertToStateUser(dbUser *models.User) *types.User {
	if dbUser == nil {
		return nil
	}

	var role types.UserRole
	switch string(dbUser.Role) {
	case "admin":
		role = types.RoleAdmin
	case "user":
		role = types.RoleUser
	default:
		role = types.RoleUser
	}

	return &types.User{
		ID:         dbUser.ID,
		TelegramID: dbUser.TelegramID,
		Username:   dbUser.Username,
		FullName:   dbUser.FullName,
		Role:       role,
	}
}
