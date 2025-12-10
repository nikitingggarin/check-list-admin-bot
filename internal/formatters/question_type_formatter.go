package formatters

import (
	"telegram-bot/internal/state_manager/types"
	"telegram-bot/internal/storage/models"
)

func FormatQuestionType(category types.QuestionCategory) string {
	switch category {
	case types.CategoryCompliance:
		return "✅ Соответствие"
	case types.CategorySingleChoice:
		return "🔘 Одиночный выбор"
	case types.CategoryMultipleChoice:
		return "☑️ Множественный выбор"
	case types.CategoryTextAnswer:
		return "📝 Текстовый ответ"
	default:
		return string(category)
	}
}

func FormatQuestionTypeModels(category models.QuestionCategory) string {
	switch category {
	case models.CategoryCompliance:
		return "✅ Соответствие"
	case models.CategorySingleChoice:
		return "🔘 Одиночный выбор"
	case models.CategoryMultipleChoice:
		return "☑️ Множественный выбор"
	case models.CategoryTextAnswer:
		return "📝 Текстовый ответ"
	default:
		return string(category)
	}
}
