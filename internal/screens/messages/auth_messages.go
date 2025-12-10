package messages

import (
	"telegram-bot/internal/state_manager/state"
)

// AuthMessage - сообщение для экрана авторизации
func AuthMessage(state *state.UserState) string {
	return "👋 Привет! Я бот для управления чек-листами.\n\n" +
		"Для начала работы необходимо авторизоваться."
}

// AdminMenuMessage - сообщение для главного меню администратора
func AdminMenuMessage(state *state.UserState) string {
	if state != nil && state.User != nil {
		return "✅ Добро пожаловать, " + state.User.Username + " " + state.User.FullName + "!\n\nВыберите действие:"
	}
	return "Главное меню администратора. Выберите действие:"
}
