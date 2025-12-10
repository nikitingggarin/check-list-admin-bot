package models

import (
	"time"
)

// AnswerOption - вариант ответа
type AnswerOption struct {
	ID         int64     `json:"id,omitempty"` // omitempty
	QuestionID int64     `json:"question_id"`
	Text       string    `json:"text"`
	IsCorrect  bool      `json:"is_correct"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
