package models

import (
	"time"
)

// Типы статусов чек-листа
type ChecklistStatus string

const (
	StatusDraft       ChecklistStatus = "draft"
	StatusPublished   ChecklistStatus = "published"
	StatusUnpublished ChecklistStatus = "unpublished"
)

// Checklist - чек-лист
type Checklist struct {
	ID        int64           `json:"id,omitempty"` // omitempty чтобы не передавать при вставке
	Name      string          `json:"name"`
	UserID    int64           `json:"user_id"`
	Status    ChecklistStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
}
