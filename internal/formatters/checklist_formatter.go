package formatters

import (
	"fmt"
	"strings"
	"telegram-bot/internal/storage/models"
)

func FormatChecklist(checklist *models.Checklist, hasBlocks bool, blocks []models.QuestionBlock, templates []models.ChecklistTemplate, questions []models.Question, answerOptions []models.AnswerOption) string {
	var sb strings.Builder

	// Заголовок в зависимости от статуса
	switch checklist.Status {
	case models.StatusPublished:
		sb.WriteString("🚀 ПРЕВЬЮ ЧЕК-ЛИСТА (ОПУБЛИКОВАН)\n")
	case models.StatusUnpublished:
		sb.WriteString("🚫 ПРЕВЬЮ ЧЕК-ЛИСТА (ОТМЕНЕН)\n")
	default: // StatusDraft
		sb.WriteString("📋 ПРЕВЬЮ ЧЕК-ЛИСТА\n")
	}

	sb.WriteString("═════════════════════\n\n")

	sb.WriteString(fmt.Sprintf("🏷️ Название: %s\n", checklist.Name))
	sb.WriteString(fmt.Sprintf("📊 ID: %d\n", checklist.ID))

	answersByQuestion := make(map[int64][]models.AnswerOption)
	for _, ao := range answerOptions {
		answersByQuestion[ao.QuestionID] = append(answersByQuestion[ao.QuestionID], ao)
	}

	if hasBlocks && len(blocks) > 0 {
		sb.WriteString("🎯 Тип: С блоками\n")
		sb.WriteString(fmt.Sprintf("🧱 Блоков: %d\n", len(blocks)))

		questionsByBlock := make(map[int64][]models.Question)
		questionMap := make(map[int64]models.Question)

		for _, q := range questions {
			questionMap[q.ID] = q
		}

		for _, t := range templates {
			if t.BlockID != nil {
				if question, exists := questionMap[t.QuestionID]; exists && question.ID != 0 {
					questionsByBlock[*t.BlockID] = append(questionsByBlock[*t.BlockID], question)
				}
			}
		}

		sb.WriteString("\n═════════════════════\n")

		for b, block := range blocks {
			sb.WriteString(fmt.Sprintf("🧱 БЛОК %d: %s\n", b+1, block.Name))
			sb.WriteString("═════════════════════\n\n")

			if blockQuestions, ok := questionsByBlock[block.ID]; ok && len(blockQuestions) > 0 {
				for i, question := range blockQuestions {
					formatQuestion(&sb, question, answersByQuestion[question.ID], i, "  ")
				}
			} else {
				sb.WriteString("  📭 В блоке нет вопросов\n\n")
			}
		}

	} else {
		sb.WriteString("🎯 Тип: Простой\n")
		sb.WriteString(fmt.Sprintf("❓ Вопросов: %d\n", len(questions)))

		sb.WriteString("\n═════════════════════\n")
		sb.WriteString("📝 ВОПРОСЫ:\n")
		sb.WriteString("═════════════════════\n\n")

		for i, question := range questions {
			formatQuestion(&sb, question, answersByQuestion[question.ID], i, "")
		}
	}

	sb.WriteString("\n═════════════════════\n")
	createdAtStr := "(дата недоступна)"
	if !checklist.CreatedAt.IsZero() && checklist.CreatedAt.Year() > 1 {
		createdAtStr = checklist.CreatedAt.Format("02.01.2006 15:04")
	}
	sb.WriteString(fmt.Sprintf("📅 Создан: %s\n", createdAtStr))
	sb.WriteString(fmt.Sprintf("📊 Статус: %s\n", checklist.Status))

	sb.WriteString("\n═════════════════════\n")
	sb.WriteString("Выберите действие:")

	return sb.String()
}
func formatQuestion(sb *strings.Builder, question models.Question, options []models.AnswerOption, index int, prefix string) {
	sb.WriteString(fmt.Sprintf("%s%d. %s\n", prefix, index+1, question.Text))
	sb.WriteString(fmt.Sprintf("%s   📌 Тип: %s\n", prefix, FormatQuestionTypeModels(question.Category)))

	if len(options) > 0 {
		sb.WriteString(fmt.Sprintf("%s   📊 Варианты ответов:\n", prefix))
		for j, option := range options {
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
