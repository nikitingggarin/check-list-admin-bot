package formatters

import (
	"fmt"
	"strings"
	"telegram-bot/internal/state_manager/types"
)

func FormatSimpleChecklistPreview(checklist *types.SimpleCheckList) string {
	var sb strings.Builder

	sb.WriteString("📋 ПРЕВЬЮ ЧЕК-ЛИСТА\n")
	sb.WriteString("═════════════════════\n\n")

	sb.WriteString("🏷️ Название: ")
	sb.WriteString(checklist.Name)
	sb.WriteString("\n")

	sb.WriteString("📊 Статус: ")
	sb.WriteString(string(checklist.Status))
	sb.WriteString("\n")

	sb.WriteString("❓ Количество вопросов: ")
	sb.WriteString(fmt.Sprintf("%d", len(checklist.Questions)))
	sb.WriteString("\n\n")

	sb.WriteString("═════════════════════\n")
	sb.WriteString("📝 ВОПРОСЫ:\n")
	sb.WriteString("═════════════════════\n\n")

	for i, question := range checklist.Questions {
		formatStateQuestion(&sb, question, i, "")
	}

	sb.WriteString("═════════════════════\n")
	sb.WriteString("Выберите действие:\n")

	return sb.String()
}

func FormatBlockedChecklistPreview(checklist *types.BlockedCheckList) string {
	var sb strings.Builder

	sb.WriteString("📋 ПРЕВЬЮ ЧЕК-ЛИСТА С БЛОКАМИ\n")
	sb.WriteString("═════════════════════\n\n")

	sb.WriteString("🏷️ Название: ")
	sb.WriteString(checklist.Name)
	sb.WriteString("\n")

	sb.WriteString("📊 Статус: ")
	sb.WriteString(string(checklist.Status))
	sb.WriteString("\n")

	sb.WriteString("🧱 Количество блоков: ")
	sb.WriteString(fmt.Sprintf("%d", len(checklist.Blocks)))
	sb.WriteString("\n\n")

	totalQuestions := 0
	for _, block := range checklist.Blocks {
		totalQuestions += len(block.Questions)
	}
	sb.WriteString("❓ Общее количество вопросов: ")
	sb.WriteString(fmt.Sprintf("%d", totalQuestions))
	sb.WriteString("\n\n")

	for b, block := range checklist.Blocks {
		sb.WriteString("═════════════════════\n")
		sb.WriteString("🧱 БЛОК ")
		sb.WriteString(fmt.Sprintf("%d", b+1))
		sb.WriteString(": ")
		sb.WriteString(block.Name)
		sb.WriteString("\n")
		sb.WriteString("═════════════════════\n\n")

		for i, question := range block.Questions {
			formatStateQuestion(&sb, question, i, "  ")
		}
	}

	sb.WriteString("═════════════════════\n")
	sb.WriteString("Выберите действие:\n")

	return sb.String()
}

func formatStateQuestion(sb *strings.Builder, question types.Question, index int, prefix string) {
	sb.WriteString(fmt.Sprintf("%s%d. %s\n", prefix, index+1, question.Text))
	sb.WriteString(fmt.Sprintf("%s   📌 Тип: %s\n", prefix, FormatQuestionType(question.Category)))

	if len(question.AnswerOptions) > 0 {
		sb.WriteString(fmt.Sprintf("%s   📊 Варианты ответов:\n", prefix))
		for j, option := range question.AnswerOptions {
			sb.WriteString(prefix + "     ")
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
