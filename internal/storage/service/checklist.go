package service

import (
	"context"
	"fmt"
	"log"
	"telegram-bot/internal/state_manager/types"
	"telegram-bot/internal/storage/models"
	"telegram-bot/internal/storage/repositories"
	"time"
)

// ChecklistService - сервис для работы с чек-листами
type ChecklistService struct {
	checklistRepo     repositories.ChecklistRepositoryInterface
	questionBlockRepo repositories.QuestionBlockRepositoryInterface
	questionRepo      repositories.QuestionRepositoryInterface
	answerOptionRepo  repositories.AnswerOptionRepositoryInterface
	templateRepo      repositories.ChecklistTemplateRepositoryInterface
	userRepo          repositories.Repositories // Добавляем репозиторий пользователей
}

// NewChecklistService создает новый сервис для чек-листов
func NewChecklistService(
	checklistRepo repositories.ChecklistRepositoryInterface,
	questionBlockRepo repositories.QuestionBlockRepositoryInterface,
	questionRepo repositories.QuestionRepositoryInterface,
	answerOptionRepo repositories.AnswerOptionRepositoryInterface,
	templateRepo repositories.ChecklistTemplateRepositoryInterface,
	userRepo repositories.Repositories, // Добавляем параметр
) *ChecklistService {
	return &ChecklistService{
		checklistRepo:     checklistRepo,
		questionBlockRepo: questionBlockRepo,
		questionRepo:      questionRepo,
		answerOptionRepo:  answerOptionRepo,
		templateRepo:      templateRepo,
		userRepo:          userRepo,
	}
}

// SaveSimpleChecklistDraft сохраняет простой чек-лист в базу данных
func (s *ChecklistService) SaveSimpleChecklistDraft(ctx context.Context, checklist *types.SimpleCheckList, telegramUserID int64) (*models.Checklist, error) {
	log.Printf("[ChecklistService] Сохранение простого чек-листа '%s' для пользователя %d", checklist.Name, telegramUserID)

	// 1. Получаем ID пользователя из базы по telegram_id
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Создаем чек-лист в базе
	dbChecklist := &models.Checklist{
		Name:      checklist.Name,
		UserID:    user.ID, // Используем ID из таблицы users
		Status:    models.StatusDraft,
		CreatedAt: time.Now(),
	}

	createdChecklist, err := s.checklistRepo.Create(ctx, dbChecklist)
	if err != nil {
		return nil, fmt.Errorf("failed to create checklist: %w", err)
	}
	log.Printf("[ChecklistService] Чек-лист создан с ID: %d", createdChecklist.ID)

	// 3. Сохраняем вопросы и варианты ответов
	for i, question := range checklist.Questions {
		// Создаем вопрос в базе (без block_id для простых чек-листов)
		dbQuestion := &models.Question{
			Text:        question.Text,
			Category:    models.QuestionCategory(question.Category),
			ChecklistID: createdChecklist.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		createdQuestion, err := s.questionRepo.Create(ctx, dbQuestion)
		if err != nil {
			return nil, fmt.Errorf("failed to create question %d: %w", i+1, err)
		}

		// 4. Сохраняем варианты ответов (если есть)
		if len(question.AnswerOptions) > 0 {
			var answerOptions []models.AnswerOption
			for _, option := range question.AnswerOptions {
				answerOptions = append(answerOptions, models.AnswerOption{
					QuestionID: createdQuestion.ID,
					Text:       option.Text,
					IsCorrect:  option.IsCorrect,
					CreatedAt:  time.Now(),
				})
			}

			if len(answerOptions) > 0 {
				_, err := s.answerOptionRepo.CreateBatch(ctx, answerOptions)
				if err != nil {
					return nil, fmt.Errorf("failed to create answer options for question %d: %w", i+1, err)
				}
			}
		}

		// 5. Создаем связку в checklist_templates (без block_id)
		template := &models.ChecklistTemplate{
			ChecklistID: createdChecklist.ID,
			QuestionID:  createdQuestion.ID,
			CreatedAt:   time.Now(),
		}

		_, err = s.templateRepo.Create(ctx, template)
		if err != nil {
			return nil, fmt.Errorf("failed to create checklist template for question %d: %w", i+1, err)
		}

		log.Printf("[ChecklistService] Вопрос %d сохранен с ID: %d", i+1, createdQuestion.ID)
	}

	log.Printf("[ChecklistService] Простой чек-лист успешно сохранен. ID: %d, вопросов: %d",
		createdChecklist.ID, len(checklist.Questions))

	return createdChecklist, nil
}

// SaveBlockedChecklistDraft сохраняет чек-лист с блоками в базу данных
func (s *ChecklistService) SaveBlockedChecklistDraft(ctx context.Context, checklist *types.BlockedCheckList, telegramUserID int64) (*models.Checklist, error) {
	log.Printf("[ChecklistService] Сохранение чек-листа с блоками '%s' для пользователя %d", checklist.Name, telegramUserID)

	// 1. Получаем ID пользователя из базы по telegram_id
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Создаем чек-лист в базе
	dbChecklist := &models.Checklist{
		Name:      checklist.Name,
		UserID:    user.ID,
		Status:    models.StatusDraft,
		CreatedAt: time.Now(),
	}

	createdChecklist, err := s.checklistRepo.Create(ctx, dbChecklist)
	if err != nil {
		return nil, fmt.Errorf("failed to create checklist: %w", err)
	}
	log.Printf("[ChecklistService] Чек-лист создан с ID: %d", createdChecklist.ID)

	// 3. Сохраняем блоки
	for blockIndex, block := range checklist.Blocks {
		// Создаем блок в базе
		dbBlock := &models.QuestionBlock{
			Name:        block.Name,
			ChecklistID: createdChecklist.ID,
			CreatedAt:   time.Now(),
		}

		createdBlock, err := s.questionBlockRepo.Create(ctx, dbBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to create block %d: %w", blockIndex+1, err)
		}
		log.Printf("[ChecklistService] Блок %d создан с ID: %d", blockIndex+1, createdBlock.ID)

		// 4. Сохраняем вопросы в блоке
		for questionIndex, question := range block.Questions {
			// Создаем вопрос в базе (БЕЗ block_id - его нет в таблице)
			dbQuestion := &models.Question{
				Text:        question.Text,
				Category:    models.QuestionCategory(question.Category),
				ChecklistID: createdChecklist.ID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			createdQuestion, err := s.questionRepo.Create(ctx, dbQuestion)
			if err != nil {
				return nil, fmt.Errorf("failed to create question %d in block %d: %w",
					questionIndex+1, blockIndex+1, err)
			}

			// 5. Сохраняем варианты ответов (если есть)
			if len(question.AnswerOptions) > 0 {
				var answerOptions []models.AnswerOption
				for _, option := range question.AnswerOptions {
					answerOptions = append(answerOptions, models.AnswerOption{
						QuestionID: createdQuestion.ID,
						Text:       option.Text,
						IsCorrect:  option.IsCorrect,
						CreatedAt:  time.Now(),
					})
				}

				if len(answerOptions) > 0 {
					_, err := s.answerOptionRepo.CreateBatch(ctx, answerOptions)
					if err != nil {
						return nil, fmt.Errorf("failed to create answer options for question %d in block %d: %w",
							questionIndex+1, blockIndex+1, err)
					}
				}
			}

			// 6. Создаем связку в checklist_templates (с указанием block_id)
			template := &models.ChecklistTemplate{
				ChecklistID: createdChecklist.ID,
				QuestionID:  createdQuestion.ID,
				BlockID:     &createdBlock.ID, // Указываем ID блока в шаблоне
				CreatedAt:   time.Now(),
			}

			_, err = s.templateRepo.Create(ctx, template)
			if err != nil {
				return nil, fmt.Errorf("failed to create checklist template for question %d in block %d: %w",
					questionIndex+1, blockIndex+1, err)
			}

			log.Printf("[ChecklistService] Вопрос %d в блоке %d сохранен с ID: %d",
				questionIndex+1, blockIndex+1, createdQuestion.ID)
		}
	}

	// Подсчитываем общее количество вопросов
	totalQuestions := 0
	for _, block := range checklist.Blocks {
		totalQuestions += len(block.Questions)
	}

	log.Printf("[ChecklistService] Чек-лист с блоками успешно сохранен. ID: %d, блоков: %d, вопросов: %d",
		createdChecklist.ID, len(checklist.Blocks), totalQuestions)

	return createdChecklist, nil
}

// PublishChecklist публикует чек-лист
func (s *ChecklistService) PublishChecklist(ctx context.Context, checklistID int64) error {
	err := s.checklistRepo.UpdateStatus(ctx, checklistID, models.StatusPublished)
	if err != nil {
		return fmt.Errorf("failed to publish checklist: %w", err)
	}

	log.Printf("[ChecklistService] Чек-лист %d опубликован", checklistID)
	return nil
}

// GetUserDrafts возвращает черновики пользователя
func (s *ChecklistService) GetUserDrafts(ctx context.Context, telegramUserID int64) ([]models.Checklist, error) {
	// Получаем ID пользователя
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	status := models.StatusDraft
	return s.checklistRepo.GetByUserID(ctx, user.ID, &status)
}

// GetUserPublished возвращает опубликованные чек-листы пользователя
func (s *ChecklistService) GetUserPublished(ctx context.Context, telegramUserID int64) ([]models.Checklist, error) {
	// Получаем ID пользователя
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	status := models.StatusPublished
	return s.checklistRepo.GetByUserID(ctx, user.ID, &status)
}

// GetChecklistByID загружает чек-лист по ID со всеми данными
func (s *ChecklistService) GetChecklistByID(ctx context.Context, checklistID int64) (*models.Checklist, []models.QuestionBlock, []models.Question, []models.AnswerOption, error) {
	// 1. Получаем чек-лист
	checklist, err := s.checklistRepo.GetByID(ctx, checklistID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get checklist: %w", err)
	}

	// 2. Получаем блоки вопросов
	blocks, err := s.questionBlockRepo.GetByChecklistID(ctx, checklistID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get question blocks: %w", err)
	}

	// 3. Получаем вопросы
	questions, err := s.questionRepo.GetByChecklistID(ctx, checklistID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get questions: %w", err)
	}

	// 4. Получаем варианты ответов
	var questionIDs []int64
	for _, q := range questions {
		questionIDs = append(questionIDs, q.ID)
	}

	var answerOptions []models.AnswerOption
	if len(questionIDs) > 0 {
		answerOptions, err = s.answerOptionRepo.GetByQuestionIDs(ctx, questionIDs)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to get answer options: %w", err)
		}
	}

	return checklist, blocks, questions, answerOptions, nil
}

// DeleteChecklist удаляет чек-лист
func (s *ChecklistService) DeleteChecklist(ctx context.Context, checklistID int64) error {
	err := s.checklistRepo.Delete(ctx, checklistID)
	if err != nil {
		return fmt.Errorf("failed to delete checklist: %w", err)
	}

	log.Printf("[ChecklistService] Чек-лист %d удален", checklistID)
	return nil
}

// GetTemplatesByChecklistID возвращает шаблоны для чек-листа
func (s *ChecklistService) GetTemplatesByChecklistID(ctx context.Context, checklistID int64) ([]models.ChecklistTemplate, error) {
	return s.templateRepo.GetByChecklistID(ctx, checklistID)
}

// UpdateChecklist обновляет чек-лист (удаляет старый и создает новый)
func (s *ChecklistService) UpdateChecklist(ctx context.Context, oldChecklistID int64, checklistData types.CheckListData, telegramUserID int64) (*models.Checklist, error) {
	log.Printf("[ChecklistService] Обновление чек-листа ID=%d", oldChecklistID)

	// 1. Удаляем старый чек-лист (каскадно удалятся все дочерние записи)
	err := s.checklistRepo.Delete(ctx, oldChecklistID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete old checklist: %w", err)
	}
	log.Printf("[ChecklistService] Старый чек-лист %d удален", oldChecklistID)

	// 2. Создаем новый чек-лист
	switch checklist := checklistData.(type) {
	case *types.SimpleCheckList:
		return s.SaveSimpleChecklistDraft(ctx, checklist, telegramUserID)
	case *types.BlockedCheckList:
		return s.SaveBlockedChecklistDraft(ctx, checklist, telegramUserID)
	default:
		return nil, fmt.Errorf("unknown checklist type")
	}
}

// GetUserUnpublished возвращает отмененные чек-листы пользователя
func (s *ChecklistService) GetUserUnpublished(ctx context.Context, telegramUserID int64) ([]models.Checklist, error) {
	// Получаем ID пользователя
	user, err := s.userRepo.GetUserByTelegramID(ctx, telegramUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	status := models.StatusUnpublished
	return s.checklistRepo.GetByUserID(ctx, user.ID, &status)
}

// UnpublishChecklist снимает чек-лист с публикации
func (s *ChecklistService) UnpublishChecklist(ctx context.Context, checklistID int64) error {
	err := s.checklistRepo.UpdateStatus(ctx, checklistID, models.StatusUnpublished)
	if err != nil {
		return fmt.Errorf("failed to unpublish checklist: %w", err)
	}

	log.Printf("[ChecklistService] Чек-лист %d снят с публикации", checklistID)
	return nil
}

// RepublishChecklist возвращает чек-лист в публикацию
func (s *ChecklistService) RepublishChecklist(ctx context.Context, checklistID int64) error {
	err := s.checklistRepo.UpdateStatus(ctx, checklistID, models.StatusPublished)
	if err != nil {
		return fmt.Errorf("failed to republish checklist: %w", err)
	}

	log.Printf("[ChecklistService] Чек-лист %d возвращен в публикацию", checklistID)
	return nil
}
