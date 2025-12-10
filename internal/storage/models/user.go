package models

import "time"

// Типизированные константы (enum-like)
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User - пользователь бота
type User struct {
	ID         int64     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	Role       UserRole  `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}
