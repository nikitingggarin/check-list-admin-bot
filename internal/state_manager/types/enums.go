package types

// UserRole - роли пользователя
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// CheckListStatus - статусы чек-листа
type CheckListStatus string

const (
	StatusDraft       CheckListStatus = "draft"
	StatusPublished   CheckListStatus = "published"
	StatusUnpublished CheckListStatus = "unpublished"
)

// QuestionCategory - категории вопросов
type QuestionCategory string

const (
	CategoryCompliance     QuestionCategory = "compliance"
	CategorySingleChoice   QuestionCategory = "single_choice"
	CategoryMultipleChoice QuestionCategory = "multiple_choice"
	CategoryTextAnswer     QuestionCategory = "text_answer"
)

// CheckListType - тип чек-листа
type CheckListType string

const (
	TypeSimple  CheckListType = "simple"
	TypeBlocked CheckListType = "blocked"
)
