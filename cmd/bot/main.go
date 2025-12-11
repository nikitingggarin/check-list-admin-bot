package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"telegram-bot/config"
	"telegram-bot/internal/handlers"
	"telegram-bot/internal/routes"

	"telegram-bot/internal/services/answers"
	"telegram-bot/internal/services/auth"
	"telegram-bot/internal/services/block_checklist"
	"telegram-bot/internal/services/menu"
	"telegram-bot/internal/services/my_checklists"
	"telegram-bot/internal/services/published_checklists"
	"telegram-bot/internal/services/question"
	"telegram-bot/internal/services/question_edit"
	"telegram-bot/internal/services/screen"
	"telegram-bot/internal/services/simple_checklist"
	"telegram-bot/internal/state_manager/manager"
	"telegram-bot/internal/storage/infrastructure"
	"telegram-bot/internal/storage/repositories"
	"telegram-bot/internal/storage/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 1. Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("🛑 Получен сигнал %v, завершаю работу...", sig)
		cancel()
		time.Sleep(2 * time.Second)
		log.Println("👋 Бот остановлен")
		os.Exit(0)
	}()

	// 2. Загрузка конфигурации
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 3. Инициализация бота
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	bot.Debug = false
	log.Printf("🤖 %s ЗАПУЩЕН", bot.Self.UserName)
	log.Println("==========================================")

	// 4. Инициализация StateManager
	stateMgr := manager.NewMemoryStateManager()

	// 5. Инициализация базы данных
	dbClient, err := infrastructure.NewDatabaseClient(cfg.Database.URL, cfg.Database.Key)
	if err != nil {
		log.Fatal("Failed to create database client:", err)
	}

	// 6. Инициализация репозиториев
	supabaseAdapter := repositories.NewSupabaseAdapter(dbClient.Client())
	userSvc := service.NewUserService(supabaseAdapter)

	// Новые репозитории для работы с чек-листами
	checklistRepo := repositories.NewChecklistRepository(dbClient.Client())
	questionBlockRepo := repositories.NewQuestionBlockRepository(dbClient.Client())
	questionRepo := repositories.NewQuestionRepository(dbClient.Client())
	answerOptionRepo := repositories.NewAnswerOptionRepository(dbClient.Client())
	templateRepo := repositories.NewChecklistTemplateRepository(dbClient.Client())

	// Инициализация сервиса чек-листов
	checklistSvc := service.NewChecklistService(
		checklistRepo,
		questionBlockRepo,
		questionRepo,
		answerOptionRepo,
		templateRepo,
		supabaseAdapter,
	)

	// Инициализация сервисов работы с бизнес логикой
	screenSvc := screen.NewScreenService(bot)
	authSvc := auth.NewAuthService(stateMgr, userSvc, screenSvc)

	// Создаем MyChecklistsService
	myChecklistsSvc := my_checklists.NewMyChecklistsService(stateMgr, screenSvc, checklistSvc)

	// Создаем PublishedChecklistsService
	publishedChecklistsSvc := published_checklists.NewPublishedChecklistsService(stateMgr, screenSvc, checklistSvc)

	// Создаем MenuService
	menuSvc := menu.NewMenuService(stateMgr, screenSvc, checklistSvc, publishedChecklistsSvc)

	// Инициализация сервисов для вопросов и ответов
	questionSvc := question.NewQuestionService(stateMgr, screenSvc)
	answersSvc := answers.NewAnswersService(stateMgr, screenSvc, questionSvc)
	questionEditSvc := question_edit.NewQuestionEditService(stateMgr, screenSvc, questionSvc, answersSvc)

	// Инициализация сервисов для чек-листов
	simpleChecklistSvc := simple_checklist.NewSimpleChecklistService(stateMgr, screenSvc, questionEditSvc, checklistSvc)
	blockChecklistSvc := block_checklist.NewBlockChecklistService(stateMgr, screenSvc, checklistSvc)

	// Инициализация роутов
	authRoute := routes.NewAuthRoutes(authSvc, screenSvc)
	menuRoute := routes.NewMenuRoutes(menuSvc, screenSvc)
	simpleChecklistRoute := routes.NewSimpleChecklistRoutes(simpleChecklistSvc, screenSvc)
	blockChecklistRoute := routes.NewBlockChecklistRoutes(blockChecklistSvc, screenSvc)
	questionRoute := routes.NewQuestionRoutes(questionSvc, screenSvc)
	answersRoute := routes.NewAnswersRoutes(answersSvc, screenSvc)
	questionEditRoute := routes.NewQuestionEditRoutes(questionEditSvc, screenSvc)
	myChecklistsRoute := routes.NewMyChecklistsRoutes(myChecklistsSvc, screenSvc)
	publishedRoute := routes.NewPublishedChecklistsRoutes(publishedChecklistsSvc, screenSvc)

	// Инициализация роутера
	router := routes.NewRouter(
		stateMgr,
		userSvc,
		authRoute,
		menuRoute,
		simpleChecklistRoute,
		blockChecklistRoute,
		questionRoute,
		answersRoute,
		questionEditRoute,
		myChecklistsRoute,
		publishedRoute,
	)

	// Инициализация обработчика обновлений
	updateHandler := handlers.NewUpdateHandler(router, stateMgr)

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	log.Println("🚀 Бот запущен и ожидает сообщений...")
	log.Println("==========================================")

	// Главный цикл с graceful shutdown
	for {
		select {
		case update := <-updates:
			go updateHandler.HandleUpdate(update)

		case <-ctx.Done():
			log.Println("🛑 Останавливаю получение новых сообщений...")
			return
		}
	}
}
