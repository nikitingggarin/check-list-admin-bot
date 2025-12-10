package configs

import (
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ScreenConfig - конфигурация экрана
type ScreenConfig struct {
	Keyboard func(state *state.UserState) tgbotapi.ReplyKeyboardMarkup
	Message  func(state *state.UserState) string
}

// ScreenConfigs - глобальная карта конфигураций экранов
var ScreenConfigs = map[string]ScreenConfig{}
