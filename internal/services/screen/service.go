package screen

import (
	"telegram-bot/internal/screens"
	"telegram-bot/internal/state_manager/state"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ScreenService отвечает за отправку экранов пользователю
type ScreenService struct {
	bot *tgbotapi.BotAPI
}

// NewScreenService создает новый сервис экранов
func NewScreenService(bot *tgbotapi.BotAPI) *ScreenService {
	return &ScreenService{bot: bot}
}

// SendScreen отправляет конкретный экран
func (s *ScreenService) SendScreen(chatID int64, screen string, userState *state.UserState) error {
	// Получаем клавиатуру и сообщение для экрана
	keyboard := screens.GetKeyboardForState(userState)
	message := screens.GetMessageForScreen(screen, userState)

	// Создаем и отправляем сообщение
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = keyboard

	_, err := s.bot.Send(msg)
	return err
}

// SendCurrentScreen отправляет текущий экран пользователя
func (s *ScreenService) SendCurrentScreen(chatID int64, userState *state.UserState) error {
	currentScreen := userState.GetCurrentScreen()
	if currentScreen == "" {
		currentScreen = "authorize-admin"
	}
	return s.SendScreen(chatID, currentScreen, userState)
}

// SendMessage отправляет простое сообщение
func (s *ScreenService) SendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := s.bot.Send(msg)
	return err
}
