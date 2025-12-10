package repositories

import (
	"context"
	"fmt"
	"strings"
	"telegram-bot/internal/storage/models"

	"github.com/nedpals/supabase-go"
)

// AnswerOptionRepository - репозиторий для работы с вариантами ответов
type AnswerOptionRepository struct {
	client *supabase.Client
}

// NewAnswerOptionRepository создает новый репозиторий для вариантов ответов
func NewAnswerOptionRepository(client *supabase.Client) *AnswerOptionRepository {
	return &AnswerOptionRepository{client: client}
}

// Create создает новый вариант ответа
func (r *AnswerOptionRepository) Create(ctx context.Context, option *models.AnswerOption) (*models.AnswerOption, error) {
	var result []models.AnswerOption

	err := r.client.DB.From("answer_options").
		Insert(option).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create answer option: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no answer option returned after creation")
	}

	return &result[0], nil
}

// CreateBatch создает несколько вариантов ответов
func (r *AnswerOptionRepository) CreateBatch(ctx context.Context, options []models.AnswerOption) ([]models.AnswerOption, error) {
	var result []models.AnswerOption

	if len(options) == 0 {
		return result, nil
	}

	err := r.client.DB.From("answer_options").
		Insert(options).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create answer options: %w", err)
	}

	return result, nil
}

// GetByQuestionID возвращает варианты ответов по ID вопроса
func (r *AnswerOptionRepository) GetByQuestionID(ctx context.Context, questionID int64) ([]models.AnswerOption, error) {
	var options []models.AnswerOption

	err := r.client.DB.From("answer_options").
		Select("*").
		Eq("question_id", fmt.Sprint(questionID)).
		Execute(&options)

	if err != nil {
		return nil, fmt.Errorf("failed to get answer options: %w", err)
	}

	return options, nil
}

// GetByQuestionIDs возвращает варианты ответов для нескольких вопросов
func (r *AnswerOptionRepository) GetByQuestionIDs(ctx context.Context, questionIDs []int64) ([]models.AnswerOption, error) {
	var options []models.AnswerOption

	if len(questionIDs) == 0 {
		return options, nil
	}

	// Создаем строку для фильтра IN
	var idStrings []string
	for _, id := range questionIDs {
		idStrings = append(idStrings, fmt.Sprint(id))
	}

	err := r.client.DB.From("answer_options").
		Select("*").
		Filter("question_id", "in", fmt.Sprintf("(%s)", strings.Join(idStrings, ","))).
		Execute(&options)

	if err != nil {
		return nil, fmt.Errorf("failed to get answer options: %w", err)
	}

	return options, nil
}
