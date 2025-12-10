package configs

import (
	"telegram-bot/internal/screens/keyboards"
	"telegram-bot/internal/screens/messages"
)

// Init регистрирует все конфигурации экранов
func Init() {
	// ========== ЭКРАНЫ АВТОРИЗАЦИИ И МЕНЮ ==========
	ScreenConfigs["authorize-admin"] = ScreenConfig{
		Keyboard: keyboards.GetAuthKeyboard,
		Message:  messages.AuthMessage,
	}

	ScreenConfigs["admin-menu"] = ScreenConfig{
		Keyboard: keyboards.GetAdminMenu,
		Message:  messages.AdminMenuMessage,
	}

	// ========== ЭКРАНЫ ДЛЯ ОПУБЛИКОВАННЫХ/ОТМЕНЕННЫХ ЧЕК-ЛИСТОВ ==========
	ScreenConfigs["published-checklists-list"] = ScreenConfig{
		Keyboard: keyboards.GetPublishedChecklistsKeyboard,
		Message:  messages.PublishedChecklistsListMessage,
	}

	ScreenConfigs["published-checklist-detail"] = ScreenConfig{
		Keyboard: keyboards.GetPublishedChecklistDetailKeyboard,
		Message:  messages.PublishedChecklistDetailMessage,
	}

	// ========== ЭКРАНЫ ДЛЯ МОИХ ЧЕК-ЛИСТОВ ==========
	ScreenConfigs["my-checklists-list"] = ScreenConfig{
		Keyboard: keyboards.GetMyChecklistsKeyboard,
		Message:  messages.MyChecklistsListMessage,
	}

	ScreenConfigs["checklist-detail"] = ScreenConfig{
		Keyboard: keyboards.GetChecklistDetailKeyboard,
		Message:  messages.ChecklistDetailMessage,
	}

	ScreenConfigs["confirm-delete-checklist"] = ScreenConfig{
		Keyboard: keyboards.GetConfirmationKeyboard,
		Message:  messages.ConfirmDeleteChecklistMessage,
	}

	// ========== ЭКРАНЫ ДЛЯ СОЗДАНИЯ ПРОСТЫХ ЧЕК-ЛИСТОВ ==========
	ScreenConfigs["create-simple-checklist-name"] = ScreenConfig{
		Keyboard: keyboards.GetBackKeyboard,
		Message:  messages.CreateSimpleChecklistNameMessage,
	}

	ScreenConfigs["simple-checklist-editor"] = ScreenConfig{
		Keyboard: keyboards.GetChecklistSimpleEditorKeyboard,
		Message:  messages.SimpleChecklistEditorMessage,
	}

	// ========== ЭКРАНЫ ДЛЯ СОЗДАНИЯ ЧЕК-ЛИСТОВ С БЛОКАМИ ==========
	ScreenConfigs["create-block-checklist-name"] = ScreenConfig{
		Keyboard: keyboards.GetBackKeyboard,
		Message:  messages.CreateBlockChecklistNameMessage,
	}

	ScreenConfigs["block-checklist-editor"] = ScreenConfig{
		Keyboard: keyboards.GetChecklistBlockEditorKeyboard,
		Message:  messages.BlockChecklistEditorMessage,
	}

	ScreenConfigs["edit-block-name"] = ScreenConfig{
		Keyboard: keyboards.GetEditBlockNameKeyboard,
		Message:  messages.EditBlockNameMessage,
	}

	ScreenConfigs["block-editor"] = ScreenConfig{
		Keyboard: keyboards.GetBlockEditorKeyboard,
		Message:  messages.BlockEditorMessage,
	}

	ScreenConfigs["confirm-exit-block-checklist"] = ScreenConfig{
		Keyboard: keyboards.GetConfirmationKeyboard,
		Message:  messages.ConfirmExitBlockChecklistMessage,
	}

	// ========== ЭКРАНЫ ДЛЯ СОЗДАНИЯ ВОПРОСОВ ==========
	ScreenConfigs["select-question-type"] = ScreenConfig{
		Keyboard: keyboards.GetQuestionTypeKeyboard,
		Message:  messages.SelectQuestionTypeMessage,
	}

	ScreenConfigs["enter-question-text"] = ScreenConfig{
		Keyboard: keyboards.GetBackKeyboard,
		Message:  messages.EnterQuestionTextMessage,
	}

	ScreenConfigs["enter-answer-options"] = ScreenConfig{
		Keyboard: keyboards.GetBackKeyboard,
		Message:  messages.EnterAnswerOptionsMessage,
	}

	ScreenConfigs["select-correct-answers"] = ScreenConfig{
		Keyboard: keyboards.GetBackKeyboard,
		Message:  messages.SelectCorrectAnswersMessage,
	}

	// ========== ЭКРАНЫ РЕДАКТИРОВАНИЯ ==========
	ScreenConfigs["edit-checklist-title"] = ScreenConfig{
		Keyboard: keyboards.GetBackKeyboard,
		Message:  messages.EditChecklistTitleMessage,
	}

	ScreenConfigs["view-question"] = ScreenConfig{
		Keyboard: keyboards.GetQuestionListKeyboard,
		Message:  messages.ViewQuestionMessage,
	}

	ScreenConfigs["edit-question-text"] = ScreenConfig{
		Keyboard: keyboards.GetBackKeyboard,
		Message:  messages.EditQuestionTextMessage,
	}

	ScreenConfigs["edit-question-type"] = ScreenConfig{
		Keyboard: keyboards.GetQuestionTypeKeyboard,
		Message:  messages.EditQuestionTypeMessage,
	}

	ScreenConfigs["confirm-delete-question"] = ScreenConfig{
		Keyboard: keyboards.GetConfirmationKeyboard,
		Message:  messages.ConfirmDeleteQuestionMessage,
	}

	// ========== ЭКРАНЫ ПРЕВЬЮ ==========
	ScreenConfigs["checklist-preview"] = ScreenConfig{
		Keyboard: keyboards.GetChecklistPreviewKeyboard,
		Message:  messages.ChecklistPreviewMessage,
	}

	// ========== ЭКРАНЫ ПОДТВЕРЖДЕНИЯ ==========
	ScreenConfigs["confirm-exit-to-main-menu"] = ScreenConfig{
		Keyboard: keyboards.GetConfirmationKeyboard,
		Message:  messages.ConfirmExitToMainMenuMessage,
	}

	// ========== ЭКРАНЫ ДЛЯ БЛОКОВ ==========
	ScreenConfigs["block-view-questions"] = ScreenConfig{
		Keyboard: keyboards.GetBlockQuestionsKeyboard,
		Message:  messages.BlockViewQuestionsMessage,
	}

	ScreenConfigs["edit-question-detail"] = ScreenConfig{
		Keyboard: keyboards.GetQuestionDetailKeyboard,
		Message:  messages.EditQuestionDetailMessage,
	}
}
