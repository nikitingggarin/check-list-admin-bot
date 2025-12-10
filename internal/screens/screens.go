package screens

import (
	"telegram-bot/internal/screens/configs"
	"telegram-bot/internal/screens/keyboards"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Инициализация конфигураций при первом импорте
func init() {
	configs.Init()
}

// GetKeyboardForState возвращает клавиатуру для текущего состояния пользователя
func GetKeyboardForState(state *state.UserState) tgbotapi.ReplyKeyboardMarkup {
	if state == nil {
		// Если состояние nil, возвращаем клавиатуру авторизации
		config, exists := configs.ScreenConfigs["authorize-admin"]
		if !exists {
			// fallback
			return keyboards.GetAuthKeyboard(nil)
		}
		return config.Keyboard(nil)
	}

	// Получаем текущий экран из breadcrumbs
	currentScreen := state.GetCurrentScreen()
	if currentScreen == "" {
		// Если нет текущего экрана, возвращаем клавиатуру авторизации
		config, exists := configs.ScreenConfigs["authorize-admin"]
		if !exists {
			return keyboards.GetAuthKeyboard(state)
		}
		return config.Keyboard(state)
	}

	// Ищем конфигурацию для текущего экрана
	config, exists := configs.ScreenConfigs[currentScreen]
	if !exists {
		// Если экран не найден в конфигах, возвращаем клавиатуру авторизации
		config, _ := configs.ScreenConfigs["authorize-admin"]
		return config.Keyboard(state)
	}

	// Возвращаем клавиатуру для текущего экрана
	return config.Keyboard(state)
}

// GetMessageForScreen возвращает сообщение для экрана
func GetMessageForScreen(screen string, state *state.UserState) string {
	config, exists := configs.ScreenConfigs[screen]
	if !exists {
		// Возвращаем сообщение для экрана авторизации
		config, exists := configs.ScreenConfigs["authorize-admin"]
		if !exists {
			return "Привет! Выберите действие:"
		}
		return config.Message(state)
	}
	return config.Message(state)
}
