package repositories

import (
	"telegram-bot/internal/storage/models"
)

// интерфейс для работы с репозиториями
type Repositories interface {
	// User методы
	GetUserByTelegramID(telegramID int64) (*models.User, error)
}

// ChecklistRepositoryInterface - интерфейс для работы с чек-листами
type ChecklistRepositoryInterface interface {
	Create(checklist *models.Checklist) (*models.Checklist, error)
	GetByID(id int64) (*models.Checklist, error)
	GetByUserID(userID int64, status *models.ChecklistStatus) ([]models.Checklist, error)
	UpdateStatus(id int64, status models.ChecklistStatus) error
	UpdateName(id int64, name string) error
	Delete(id int64) error
}

// QuestionBlockRepositoryInterface - интерфейс для работы с блоками вопросов
type QuestionBlockRepositoryInterface interface {
	Create(block *models.QuestionBlock) (*models.QuestionBlock, error)
	CreateBatch(blocks []models.QuestionBlock) ([]models.QuestionBlock, error)
	GetByChecklistID(checklistID int64) ([]models.QuestionBlock, error)
	UpdateName(id int64, name string) error
}

// QuestionRepositoryInterface - интерфейс для работы с вопросами
type QuestionRepositoryInterface interface {
	Create(question *models.Question) (*models.Question, error)
	CreateBatch(questions []models.Question) ([]models.Question, error)
	GetByChecklistID(checklistID int64) ([]models.Question, error)
	UpdateText(id int64, text string) error
	UpdateCategory(id int64, category models.QuestionCategory) error
}

// AnswerOptionRepositoryInterface - интерфейс для работы с вариантами ответов
type AnswerOptionRepositoryInterface interface {
	Create(option *models.AnswerOption) (*models.AnswerOption, error)
	CreateBatch(options []models.AnswerOption) ([]models.AnswerOption, error)
	GetByQuestionID(questionID int64) ([]models.AnswerOption, error)
	GetByQuestionIDs(questionIDs []int64) ([]models.AnswerOption, error)
}

// ChecklistTemplateRepositoryInterface - интерфейс для работы с шаблонами
type ChecklistTemplateRepositoryInterface interface {
	Create(template *models.ChecklistTemplate) (*models.ChecklistTemplate, error)
	CreateBatch(templates []models.ChecklistTemplate) ([]models.ChecklistTemplate, error)
	GetByChecklistID(checklistID int64) ([]models.ChecklistTemplate, error)
}
