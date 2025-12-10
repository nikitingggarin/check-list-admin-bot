package models

import (
	"time"
)

// ChecklistTemplate - шаблон связки чек-листа, вопроса и блока
type ChecklistTemplate struct {
	ID          int64     `json:"id,omitempty"` // omitempty
	ChecklistID int64     `json:"checklist_id"`
	QuestionID  int64     `json:"question_id"`
	BlockID     *int64    `json:"block_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
