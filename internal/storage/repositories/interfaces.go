package repositories

import (
	"context"
	"telegram-bot/internal/storage/models"
)

// интерфейс для работы с репозиториями
type Repositories interface {
	// User методы
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
}

// ChecklistRepositoryInterface - интерфейс для работы с чек-листами
type ChecklistRepositoryInterface interface {
	Create(ctx context.Context, checklist *models.Checklist) (*models.Checklist, error)
	GetByID(ctx context.Context, id int64) (*models.Checklist, error)
	GetByUserID(ctx context.Context, userID int64, status *models.ChecklistStatus) ([]models.Checklist, error)
	UpdateStatus(ctx context.Context, id int64, status models.ChecklistStatus) error
	UpdateName(ctx context.Context, id int64, name string) error
	Delete(ctx context.Context, id int64) error
}

// QuestionBlockRepositoryInterface - интерфейс для работы с блоками вопросов
type QuestionBlockRepositoryInterface interface {
	Create(ctx context.Context, block *models.QuestionBlock) (*models.QuestionBlock, error)
	CreateBatch(ctx context.Context, blocks []models.QuestionBlock) ([]models.QuestionBlock, error)
	GetByChecklistID(ctx context.Context, checklistID int64) ([]models.QuestionBlock, error)
	UpdateName(ctx context.Context, id int64, name string) error
}

// QuestionRepositoryInterface - интерфейс для работы с вопросами
type QuestionRepositoryInterface interface {
	Create(ctx context.Context, question *models.Question) (*models.Question, error)
	CreateBatch(ctx context.Context, questions []models.Question) ([]models.Question, error)
	GetByChecklistID(ctx context.Context, checklistID int64) ([]models.Question, error)
	UpdateText(ctx context.Context, id int64, text string) error
	UpdateCategory(ctx context.Context, id int64, category models.QuestionCategory) error
}

// AnswerOptionRepositoryInterface - интерфейс для работы с вариантами ответов
type AnswerOptionRepositoryInterface interface {
	Create(ctx context.Context, option *models.AnswerOption) (*models.AnswerOption, error)
	CreateBatch(ctx context.Context, options []models.AnswerOption) ([]models.AnswerOption, error)
	GetByQuestionID(ctx context.Context, questionID int64) ([]models.AnswerOption, error)
	GetByQuestionIDs(ctx context.Context, questionIDs []int64) ([]models.AnswerOption, error)
}

// ChecklistTemplateRepositoryInterface - интерфейс для работы с шаблонами
type ChecklistTemplateRepositoryInterface interface {
	Create(ctx context.Context, template *models.ChecklistTemplate) (*models.ChecklistTemplate, error)
	CreateBatch(ctx context.Context, templates []models.ChecklistTemplate) ([]models.ChecklistTemplate, error)
	GetByChecklistID(ctx context.Context, checklistID int64) ([]models.ChecklistTemplate, error)
}
