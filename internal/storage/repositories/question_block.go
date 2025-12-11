package repositories

import (
	"fmt"
	"telegram-bot/internal/storage/models"

	"github.com/nedpals/supabase-go"
)

// QuestionBlockRepository - репозиторий для работы с блоками вопросов
type QuestionBlockRepository struct {
	client *supabase.Client
}

// NewQuestionBlockRepository создает новый репозиторий для блоков вопросов
func NewQuestionBlockRepository(client *supabase.Client) *QuestionBlockRepository {
	return &QuestionBlockRepository{client: client}
}

// Create создает новый блок вопросов
func (r *QuestionBlockRepository) Create(block *models.QuestionBlock) (*models.QuestionBlock, error) {
	var result []models.QuestionBlock

	err := r.client.DB.From("question_blocks").
		Insert(block).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create question block: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no question block returned after creation")
	}

	return &result[0], nil
}

// CreateBatch создает несколько блоков вопросов
func (r *QuestionBlockRepository) CreateBatch(blocks []models.QuestionBlock) ([]models.QuestionBlock, error) {
	var result []models.QuestionBlock

	if len(blocks) == 0 {
		return result, nil
	}

	err := r.client.DB.From("question_blocks").
		Insert(blocks).
		Execute(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create question blocks: %w", err)
	}

	return result, nil
}

// GetByChecklistID возвращает блоки вопросов по ID чек-листа
func (r *QuestionBlockRepository) GetByChecklistID(checklistID int64) ([]models.QuestionBlock, error) {
	var blocks []models.QuestionBlock

	err := r.client.DB.From("question_blocks").
		Select("*").
		Eq("checklist_id", fmt.Sprint(checklistID)).
		Execute(&blocks)

	if err != nil {
		return nil, fmt.Errorf("failed to get question blocks: %w", err)
	}

	return blocks, nil
}

// UpdateName обновляет название блока
func (r *QuestionBlockRepository) UpdateName(id int64, name string) error {
	data := map[string]interface{}{
		"name": name,
	}

	var result []models.QuestionBlock
	err := r.client.DB.From("question_blocks").
		Update(data).
		Eq("id", fmt.Sprint(id)).
		Execute(&result)

	if err != nil {
		return fmt.Errorf("failed to update question block name: %w", err)
	}

	return nil
}
