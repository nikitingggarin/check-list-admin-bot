package buttons

// Невидимый символ Zero Width Space (U+200B)
const zeroWidthSpace = ""

// Константы для всех кнопок
const (
	// Админ меню
	BtnAuth                  = "🔐 Авторизоваться" + zeroWidthSpace
	BtnCreateSimpleChecklist = "📋 Создать чек-лист" + zeroWidthSpace
	BtnCreateBlockChecklist  = "🧱 Создать чек-лист с блоками" + zeroWidthSpace
	BtnMyChecklists          = "📝 Мои чек-листы" + zeroWidthSpace
	BtnPublished             = "🚀 Опубликованные" + zeroWidthSpace
	BtnCanceled              = "🚫 Отмененные публикации" + zeroWidthSpace
	BtnStatistics            = "📊 Статистика" + zeroWidthSpace
	BtnLogout                = "🚪 Выход" + zeroWidthSpace

	// Типы вопросов
	BtnCompliance     = "✅ Соответствие" + zeroWidthSpace
	BtnSingleChoice   = "🔘 Одиночный выбор" + zeroWidthSpace
	BtnMultipleChoice = "☑️ Множественный выбор" + zeroWidthSpace
	BtnTextAnswer     = "📝 Текстовый ответ" + zeroWidthSpace

	// Навигация
	BtnBack            = "◀️ Назад" + zeroWidthSpace
	BtnBackToMainMenu  = "◀️ Назад в главное меню" + zeroWidthSpace
	BtnBackToBlockList = "◀️ К списку блоков" + zeroWidthSpace

	// Редактирование вопросов
	BtnEditQuestionText = "✏️ Изменить текст вопроса" + zeroWidthSpace
	BtnEditQuestionType = "✏️ Изменить тип вопроса" + zeroWidthSpace
	BtnDeleteQuestion   = "🗑️ Удалить вопрос" + zeroWidthSpace
	BtnEdit             = "✏️ Редактировать" + zeroWidthSpace

	// Подтверждение
	BtnYes = "🟢 Да" + zeroWidthSpace
	BtnNo  = "🔴 Нет" + zeroWidthSpace

	// Редактор чек-листа
	BtnAddQuestion             = "➕ Добавить вопрос" + zeroWidthSpace
	BtnAddBlock                = "➕ Добавить блок" + zeroWidthSpace
	BtnEditTitleBlockChecklist = "✏️ Редактировать название блока" + zeroWidthSpace
	BtnPreview                 = "👁️ Посмотреть превью" + zeroWidthSpace
	BtnEditTitleChecklist      = "✏️ Редактировать название" + zeroWidthSpace
	BtnEditQuestionChecklist   = "✏️ Редактировать вопросы" + zeroWidthSpace
	BtnPreviewBlock            = "👁️ Превью блока" + zeroWidthSpace
	BtnPublishChecklist        = "🚀 Опубликовать" + zeroWidthSpace

	// BtnEditTitleBlockChecklist = "✏️ Название блока"

	BtnSaveDraft   = "💾 Сохранить черновик" + zeroWidthSpace
	BtnSavePublish = "💾🚀 Сохранить и опубликовать" + zeroWidthSpace

	BtnUnPublish = "🚫 Снять с публикации"
	BtnPublish   = "🚀 Вернуть в публикацию"

	BtnDeleteCheckList = "🗑️ Удалить"
	BtnEditCheckList   = "✏️ Редактировать чек-лист" + zeroWidthSpace
)
