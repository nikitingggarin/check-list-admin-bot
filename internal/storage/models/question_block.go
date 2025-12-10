package models

import (
	"time"
)

// QuestionBlock - блок вопросов
type QuestionBlock struct {
	ID          int64     `json:"id,omitempty"` // omitempty
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ChecklistID int64     `json:"checklist_id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
