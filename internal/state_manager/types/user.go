package types

// User - структура пользователя
type User struct {
	ID         int64    `json:"id"`
	TelegramID int64    `json:"telegram_id"`
	Username   string   `json:"username"`
	FullName   string   `json:"full_name"`
	Role       UserRole `json:"role"`
}
