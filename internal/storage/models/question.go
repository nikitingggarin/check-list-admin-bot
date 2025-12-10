package models

import (
	"time"
)

// Типы вопросов
type QuestionCategory string

const (
	CategoryCompliance     QuestionCategory = "compliance"
	CategorySingleChoice   QuestionCategory = "single_choice"
	CategoryMultipleChoice QuestionCategory = "multiple_choice"
	CategoryTextAnswer     QuestionCategory = "text_answer"
)

// Question - вопрос
type Question struct {
	ID          int64            `json:"id,omitempty"` // omitempty
	Text        string           `json:"text"`
	Category    QuestionCategory `json:"category"`
	ChecklistID int64            `json:"checklist_id"`
	CreatedAt   time.Time        `json:"created_at,omitempty"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty"`
}
