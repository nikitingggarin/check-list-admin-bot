package repositories

import (
	"context"
	"fmt"
	"telegram-bot/internal/storage/models"

	"github.com/nedpals/supabase-go"
)

// QuestionRepository - репозиторий для работы с вопросами
type QuestionRepository struct {
	client *supabase.Client
}

// NewQuestionRepository создает новый репозиторий для вопросов
func NewQuestionRepository(client *supabase.Client) *QuestionRepository {
	return &QuestionRepository{client: client}
}

// Create создает новый вопрос
func (r *QuestionRepository) Create(ctx context.Context, question *models.Question) (*models.Question, error) {
	var result []models.Question

	err := r.client.DB.From("questions").
		Insert(question).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create question: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no question returned after creation")
	}

	return &result[0], nil
}

// CreateBatch создает несколько вопросов
func (r *QuestionRepository) CreateBatch(ctx context.Context, questions []models.Question) ([]models.Question, error) {
	var result []models.Question

	if len(questions) == 0 {
		return result, nil
	}

	err := r.client.DB.From("questions").
		Insert(questions).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create questions: %w", err)
	}

	return result, nil
}

// GetByChecklistID возвращает вопросы по ID чек-листа
func (r *QuestionRepository) GetByChecklistID(ctx context.Context, checklistID int64) ([]models.Question, error) {
	var questions []models.Question

	err := r.client.DB.From("questions").
		Select("*").
		Eq("checklist_id", fmt.Sprint(checklistID)).
		Execute(&questions)

	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	return questions, nil
}

// UpdateText обновляет текст вопроса
func (r *QuestionRepository) UpdateText(ctx context.Context, id int64, text string) error {
	data := map[string]interface{}{
		"text": text,
	}

	var result []models.Question
	err := r.client.DB.From("questions").
		Update(data).
		Eq("id", fmt.Sprint(id)).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to update question text: %w", err)
	}

	return nil
}

// UpdateCategory обновляет тип вопроса
func (r *QuestionRepository) UpdateCategory(ctx context.Context, id int64, category models.QuestionCategory) error {
	data := map[string]interface{}{
		"category": category,
	}

	var result []models.Question
	err := r.client.DB.From("questions").
		Update(data).
		Eq("id", fmt.Sprint(id)).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to update question category: %w", err)
	}

	return nil
}
