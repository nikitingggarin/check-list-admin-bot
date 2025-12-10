package repositories

import (
	"context"
	"fmt"
	"log"
	"telegram-bot/internal/storage/models"

	"github.com/nedpals/supabase-go"
)

// ChecklistRepository - репозиторий для работы с чек-листами
type ChecklistRepository struct {
	client *supabase.Client
}

// NewChecklistRepository создает новый репозиторий для чек-листов
func NewChecklistRepository(client *supabase.Client) *ChecklistRepository {
	return &ChecklistRepository{client: client}
}

// Create создает новый чек-лист
func (r *ChecklistRepository) Create(ctx context.Context, checklist *models.Checklist) (*models.Checklist, error) {
	var result []models.Checklist

	err := r.client.DB.From("checklists").
		Insert(checklist).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create checklist: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no checklist returned after creation")
	}

	return &result[0], nil
}

// GetByID возвращает чек-лист по ID
func (r *ChecklistRepository) GetByID(ctx context.Context, id int64) (*models.Checklist, error) {
	var checklists []models.Checklist

	err := r.client.DB.From("checklists").
		Select("*").
		Eq("id", fmt.Sprint(id)).
		Execute(&checklists)

	if err != nil {
		return nil, fmt.Errorf("failed to get checklist: %w", err)
	}

	if len(checklists) == 0 {
		return nil, fmt.Errorf("checklist not found: %d", id)
	}

	return &checklists[0], nil
}

// GetByUserID возвращает чек-листы пользователя
func (r *ChecklistRepository) GetByUserID(ctx context.Context, userID int64, status *models.ChecklistStatus) ([]models.Checklist, error) {
	var checklists []models.Checklist

	query := r.client.DB.From("checklists").
		Select("*").
		Eq("user_id", fmt.Sprint(userID))

	if status != nil {
		query = query.Eq("status", string(*status))
	}

	err := query.Execute(&checklists)

	if err != nil {
		return nil, fmt.Errorf("failed to get user checklists: %w", err)
	}

	return checklists, nil
}

// UpdateStatus обновляет статус чек-листа
func (r *ChecklistRepository) UpdateStatus(ctx context.Context, id int64, status models.ChecklistStatus) error {
	data := map[string]interface{}{
		"status": status,
	}

	var result []models.Checklist
	err := r.client.DB.From("checklists").
		Update(data).
		Eq("id", fmt.Sprint(id)).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to update checklist status: %w", err)
	}

	return nil
}

// UpdateName обновляет название чек-листа
func (r *ChecklistRepository) UpdateName(ctx context.Context, id int64, name string) error {
	data := map[string]interface{}{
		"name": name,
	}

	var result []models.Checklist
	err := r.client.DB.From("checklists").
		Update(data).
		Eq("id", fmt.Sprint(id)).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to update checklist name: %w", err)
	}

	return nil
}

// Delete удаляет чек-лист по ID (с каскадным удалением)
func (r *ChecklistRepository) Delete(ctx context.Context, id int64) error {
	var result []models.Checklist

	err := r.client.DB.From("checklists").
		Delete().
		Eq("id", fmt.Sprint(id)).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to delete checklist: %w", err)
	}

	log.Printf("[ChecklistRepository] Чек-лист %d удален", id)
	return nil
}
