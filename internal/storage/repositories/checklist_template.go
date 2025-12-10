package repositories

import (
	"context"
	"fmt"
	"telegram-bot/internal/storage/models"

	"github.com/nedpals/supabase-go"
)

// ChecklistTemplateRepository - репозиторий для работы с шаблонами чек-листов
type ChecklistTemplateRepository struct {
	client *supabase.Client
}

// NewChecklistTemplateRepository создает новый репозиторий для шаблонов
func NewChecklistTemplateRepository(client *supabase.Client) *ChecklistTemplateRepository {
	return &ChecklistTemplateRepository{client: client}
}

// Create создает новую связку шаблона
func (r *ChecklistTemplateRepository) Create(ctx context.Context, template *models.ChecklistTemplate) (*models.ChecklistTemplate, error) {
	var result []models.ChecklistTemplate

	err := r.client.DB.From("checklist_templates").
		Insert(template).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create checklist template: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no checklist template returned after creation")
	}

	return &result[0], nil
}

// CreateBatch создает несколько связок шаблонов
func (r *ChecklistTemplateRepository) CreateBatch(ctx context.Context, templates []models.ChecklistTemplate) ([]models.ChecklistTemplate, error) {
	var result []models.ChecklistTemplate

	if len(templates) == 0 {
		return result, nil
	}

	err := r.client.DB.From("checklist_templates").
		Insert(templates).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create checklist templates: %w", err)
	}

	return result, nil
}

// GetByChecklistID возвращает шаблоны по ID чек-листа
func (r *ChecklistTemplateRepository) GetByChecklistID(ctx context.Context, checklistID int64) ([]models.ChecklistTemplate, error) {
	var templates []models.ChecklistTemplate

	err := r.client.DB.From("checklist_templates").
		Select("*").
		Eq("checklist_id", fmt.Sprint(checklistID)).
		Execute(&templates)

	if err != nil {
		return nil, fmt.Errorf("failed to get checklist templates: %w", err)
	}

	return templates, nil
}
