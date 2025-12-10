package debug

import (
	"fmt"
	"strings"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
)

// DebugState возвращает состояние в формате структуры (как в JS объекте)
func DebugState(state *state.UserState) string {
	if state == nil {
		return "Состояние: nil"
	}

	var sb strings.Builder

	// User
	sb.WriteString("{\n")
	sb.WriteString(" user: {\n")
	if state.User != nil {
		sb.WriteString(fmt.Sprintf("             id: %d,\n", state.User.ID))
		sb.WriteString(fmt.Sprintf("             telegram_id: %d,\n", state.User.TelegramID))
		sb.WriteString(fmt.Sprintf("             username: %s,\n", state.User.Username))
		sb.WriteString(fmt.Sprintf("             full_name: %s,\n", state.User.FullName))
		sb.WriteString(fmt.Sprintf("             role: %s\n", state.User.Role))
	} else {
		sb.WriteString("             id: null,\n")
		sb.WriteString("             telegram_id: null,\n")
		sb.WriteString("             username: null,\n")
		sb.WriteString("             full_name: null,\n")
		sb.WriteString("             role: null\n")
	}
	sb.WriteString("           },\n")

	// CurrentScreen (ЗАМЕНЯЕМ BreadCrumbs)
	sb.WriteString(fmt.Sprintf(" currentScreen: \"%s\",\n", state.CurrentScreen))

	// CurrentCheckList
	sb.WriteString(" currentCheckList: ")
	if state.CurrentCheckList == nil {
		sb.WriteString("null\n")
	} else {
		switch cl := state.CurrentCheckList.(type) {
		case *types.SimpleCheckList:
			debugSimpleCheckList(&sb, cl)
		case *types.BlockedCheckList:
			debugBlockedCheckList(&sb, cl)
		default:
			sb.WriteString("unknown type\n")
		}
	}

	sb.WriteString("}")
	return sb.String()
}

func debugSimpleCheckList(sb *strings.Builder, checkList *types.SimpleCheckList) {
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("        id: %d,\n", checkList.ID))
	sb.WriteString(fmt.Sprintf("        name: \"%s\",\n", checkList.Name))
	sb.WriteString(fmt.Sprintf("        status: %s,\n", checkList.Status))

	// Questions
	sb.WriteString("        questions: [\n")
	for i, q := range checkList.Questions {
		if i > 0 {
			sb.WriteString(",\n")
		}
		sb.WriteString("          {\n")
		sb.WriteString(fmt.Sprintf("            text: \"%s\",\n", q.Text))
		sb.WriteString(fmt.Sprintf("            category: %s,\n", q.Category))

		// Answer options
		if len(q.AnswerOptions) > 0 {
			sb.WriteString("            answer_options: [\n")
			for j, opt := range q.AnswerOptions {
				if j > 0 {
					sb.WriteString(",\n")
				}
				sb.WriteString("              {\n")
				sb.WriteString(fmt.Sprintf("                text: \"%s\",\n", opt.Text))
				sb.WriteString(fmt.Sprintf("                is_correct: %v\n", opt.IsCorrect))
				sb.WriteString("              }")
			}
			sb.WriteString("\n            ]\n")
		} else {
			sb.WriteString("            answer_options: []\n")
		}
		sb.WriteString("          }")
	}
	if len(checkList.Questions) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("        ]\n")
	sb.WriteString("    } \n")
}

func debugBlockedCheckList(sb *strings.Builder, checkList *types.BlockedCheckList) {
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("        id: %d,\n", checkList.ID))
	sb.WriteString(fmt.Sprintf("        name: \"%s\",\n", checkList.Name))
	sb.WriteString(fmt.Sprintf("        status: %s,\n", checkList.Status))

	// Blocks
	sb.WriteString("        blocks: [\n")
	for i, block := range checkList.Blocks {
		if i > 0 {
			sb.WriteString(",\n")
		}
		sb.WriteString("          {\n")
		sb.WriteString(fmt.Sprintf("            name: \"%s\",\n", block.Name))

		// Questions in block
		sb.WriteString("            questions: [\n")
		for j, q := range block.Questions {
			if j > 0 {
				sb.WriteString(",\n")
			}
			sb.WriteString("              {\n")
			sb.WriteString(fmt.Sprintf("                text: \"%s\",\n", q.Text))
			sb.WriteString(fmt.Sprintf("                category: %s,\n", q.Category))

			// Answer options
			if len(q.AnswerOptions) > 0 {
				sb.WriteString("                answer_options: [\n")
				for k, opt := range q.AnswerOptions {
					if k > 0 {
						sb.WriteString(",\n")
					}
					sb.WriteString("                  {\n")
					sb.WriteString(fmt.Sprintf("                    text: \"%s\",\n", opt.Text))
					sb.WriteString(fmt.Sprintf("                    is_correct: %v\n", opt.IsCorrect))
					sb.WriteString("                  }")
				}
				sb.WriteString("\n                ]\n")
			} else {
				sb.WriteString("                answer_options: []\n")
			}
			sb.WriteString("              }")
		}
		if len(block.Questions) > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("            ]\n")
		sb.WriteString("          }")
	}
	if len(checkList.Blocks) > 0 {
		sb.WriteString("\n")
	}
	sb.WriteString("        ]\n")
	sb.WriteString("    }\n")
}
