package my_checklists

import (
	"log"
	"telegram-bot/internal/state_manager/state"
	"telegram-bot/internal/state_manager/types"
	"telegram-bot/internal/storage/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *MyChecklistsService) handleEditSimpleChecklist(userID int64, update tgbotapi.Update, userState *state.UserState,
	checklist *models.Checklist, questions []models.Question, answerOptions []models.AnswerOption) {

	// Создаем чек-лист в памяти из данных базы
	simpleChecklist := &types.SimpleCheckList{
		ID:        checklist.ID,
		Name:      checklist.Name,
		Status:    types.CheckListStatus(checklist.Status),
		Questions: make([]types.Question, 0),
	}

	// Группируем ответы по вопросам
	answersByQuestion := make(map[int64][]types.AnswerOption)
	for _, ao := range answerOptions {
		answersByQuestion[ao.QuestionID] = append(answersByQuestion[ao.QuestionID],
			types.NewAnswerOption(ao.Text, ao.IsCorrect))
	}

	// Добавляем вопросы
	for _, q := range questions {
		question := types.NewQuestion(q.Text, types.QuestionCategory(q.Category))
		if opts, ok := answersByQuestion[q.ID]; ok {
			question.AnswerOptions = opts
		}
		simpleChecklist.Questions = append(simpleChecklist.Questions, question)
	}

	// Устанавливаем чек-лист в состояние
	userState.SetSimpleCheckList(simpleChecklist)
	s.stateMgr.SetState(userID, userState)

	// Переходим в редактор
	s.stateMgr.NavigateTo(userID, "simple-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[MyChecklistsService] ✅ Начато редактирование простого чек-листа ID=%d, вопросов=%d",
		checklist.ID, len(questions))
}

func (s *MyChecklistsService) handleEditBlockChecklist(userID int64, update tgbotapi.Update, userState *state.UserState,
	checklist *models.Checklist, blocks []models.QuestionBlock, questions []models.Question, answerOptions []models.AnswerOption) {

	// Создаем чек-лист с блоками в памяти
	blockedChecklist := &types.BlockedCheckList{
		ID:     checklist.ID,
		Name:   checklist.Name,
		Status: types.CheckListStatus(checklist.Status),
		Blocks: make([]types.Block, 0),
	}

	// Получаем шаблоны для группировки вопросов по блокам
	templates, err := s.checklistSvc.GetTemplatesByChecklistID(checklist.ID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "❌ Ошибка при загрузке структуры чек-листа: "+err.Error())
		return
	}

	// Группируем вопросы по блокам
	questionsByBlock := make(map[int64][]models.Question)
	questionsMap := make(map[int64]models.Question)

	for _, q := range questions {
		questionsMap[q.ID] = q
	}

	for _, t := range templates {
		if t.BlockID != nil && questionsMap[t.QuestionID].ID != 0 {
			questionsByBlock[*t.BlockID] = append(questionsByBlock[*t.BlockID], questionsMap[t.QuestionID])
		}
	}

	// Группируем ответы по вопросам
	answersByQuestion := make(map[int64][]types.AnswerOption)
	for _, ao := range answerOptions {
		answersByQuestion[ao.QuestionID] = append(answersByQuestion[ao.QuestionID],
			types.NewAnswerOption(ao.Text, ao.IsCorrect))
	}

	// Создаем блоки
	for _, block := range blocks {
		typesBlock := types.NewBlock(block.Name)

		// Добавляем вопросы в блок
		if blockQuestions, ok := questionsByBlock[block.ID]; ok {
			for _, q := range blockQuestions {
				question := types.NewQuestion(q.Text, types.QuestionCategory(q.Category))
				if opts, ok := answersByQuestion[q.ID]; ok {
					question.AnswerOptions = opts
				}
				typesBlock.AddQuestion(question)
			}
		}

		blockedChecklist.AddBlock(typesBlock)
	}

	// Устанавливаем чек-лист в состояние
	userState.SetBlockedCheckList(blockedChecklist)
	s.stateMgr.SetState(userID, userState)

	// Переходим в редактор
	s.stateMgr.NavigateTo(userID, "block-checklist-editor")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[MyChecklistsService] ✅ Начато редактирование чек-листа с блоками ID=%d, блоков=%d",
		checklist.ID, len(blocks))
}

// refreshChecklistsList обновляет список чек-листов и показывает обновленную клавиатуру
func (s *MyChecklistsService) refreshChecklistsList(userID int64, update tgbotapi.Update, userState *state.UserState) {

	// Загружаем обновленный список черновиков
	drafts, err := s.checklistSvc.GetUserDrafts(userID)
	if err != nil {
		s.screenSvc.SendMessage(update.Message.Chat.ID, "⚠️ Чек-лист удален, но не удалось обновить список: "+err.Error())
		// Все равно возвращаемся к списку
		s.stateMgr.NavigateTo(userID, "my-checklists-list")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
		return
	}

	if len(drafts) == 0 {
		// Если черновиков не осталось
		s.screenSvc.SendMessage(update.Message.Chat.ID, "📭 У вас больше нет черновиков чек-листов.")

		// Очищаем список чек-листов из состояния
		delete(userState.Data, "my_checklists")
		s.stateMgr.SetState(userID, userState)

		// Возвращаемся в главное меню
		s.stateMgr.NavigateTo(userID, "admin-menu")
		s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
		return
	}

	// Сохраняем обновленный список в состояние
	userState.Data["my_checklists"] = drafts
	s.stateMgr.SetState(userID, userState)

	// Переходим на экран списка чек-листов
	s.stateMgr.NavigateTo(userID, "my-checklists-list")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)

	log.Printf("[MyChecklistsService] ✅ Список чек-листов обновлен, осталось: %d", len(drafts))
}

// refreshChecklistsListOnBack обновляет список при возврате из деталей
func (s *MyChecklistsService) refreshChecklistsListOnBack(userID int64, update tgbotapi.Update, userState *state.UserState) {

	// Загружаем обновленный список черновиков
	drafts, err := s.checklistSvc.GetUserDrafts(userID)
	if err != nil {
		log.Printf("[MyChecklistsService] Ошибка при обновлении списка: %v", err)
		// Продолжаем с текущим списком
	} else {
		// Сохраняем обновленный список
		userState.Data["my_checklists"] = drafts
		s.stateMgr.SetState(userID, userState)
	}

	// Очищаем данные текущего чек-листа
	delete(userState.Data, "current_checklist")
	delete(userState.Data, "has_blocks")
	delete(userState.Data, "total_questions")
	delete(userState.Data, "checklist_blocks")
	delete(userState.Data, "checklist_questions")
	delete(userState.Data, "checklist_answer_options")

	s.stateMgr.SetState(userID, userState)
	s.stateMgr.NavigateTo(userID, "my-checklists-list")
	s.screenSvc.SendCurrentScreen(update.Message.Chat.ID, userState)
}
