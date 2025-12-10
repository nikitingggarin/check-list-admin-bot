package formatters

import (
	"fmt"
	"strings"
	"telegram-bot/internal/state_manager/types"
)

func FormatBlockPreview(block types.Block) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("🧱 ПРЕВЬЮ БЛОКА: %s\n", block.Name))
	sb.WriteString("═════════════════════\n\n")

	sb.WriteString(fmt.Sprintf("📊 Количество вопросов: %d\n\n", len(block.Questions)))

	sb.WriteString("═════════════════════\n")
	sb.WriteString("📝 ВОПРОСЫ:\n")
	sb.WriteString("═════════════════════\n\n")

	for i, question := range block.Questions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, question.Text))
		sb.WriteString(fmt.Sprintf("   📌 Тип: %s\n", FormatQuestionType(question.Category)))

		if len(question.AnswerOptions) > 0 {
			sb.WriteString("   📊 Варианты ответов:\n")
			for j, option := range question.AnswerOptions {
				sb.WriteString("     ")
				if option.IsCorrect {
					sb.WriteString("✅ ")
				} else {
					sb.WriteString("   ")
				}
				sb.WriteString(fmt.Sprintf("%d. %s\n", j+1, option.Text))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("═════════════════════\n")
	sb.WriteString("Вы остаетесь в редакторе блока\n")

	return sb.String()
}
