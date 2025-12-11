package repositories

import (
	"fmt"
	"telegram-bot/internal/storage/models"

	"github.com/nedpals/supabase-go"
)

type SupabaseAdapter struct {
	client *supabase.Client
}

func NewSupabaseAdapter(client *supabase.Client) *SupabaseAdapter {
	return &SupabaseAdapter{client: client}
}

func (a *SupabaseAdapter) GetUserByTelegramID(telegramID int64) (*models.User, error) {
	var users []models.User

	err := a.client.DB.From("users").
		Select("*").
		Eq("telegram_id", fmt.Sprint(telegramID)).
		Execute(&users)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found: %d", telegramID)
	}
	return &users[0], nil
}
