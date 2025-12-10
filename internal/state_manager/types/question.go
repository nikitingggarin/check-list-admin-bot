package types

// AnswerOption - вариант ответа
type AnswerOption struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// NewAnswerOption создает вариант ответа
func NewAnswerOption(text string, isCorrect bool) AnswerOption {
	return AnswerOption{
		Text:      text,
		IsCorrect: isCorrect,
	}
}

// Question - вопрос
type Question struct {
	Text          string           `json:"text"`
	Category      QuestionCategory `json:"category"`
	AnswerOptions []AnswerOption   `json:"answer_options,omitempty"`
}

// NewQuestion создает новый вопрос
func NewQuestion(text string, category QuestionCategory) Question {
	return Question{
		Text:          text,
		Category:      category,
		AnswerOptions: make([]AnswerOption, 0),
	}
}

// AddAnswerOption добавляет вариант ответа к вопросу
func (q *Question) AddAnswerOption(option AnswerOption) {
	q.AnswerOptions = append(q.AnswerOptions, option)
}

// Block - блок вопросов
type Block struct {
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
}

// NewBlock создает новый блок
func NewBlock(name string) Block {
	return Block{
		Name:      name,
		Questions: make([]Question, 0),
	}
}

// AddQuestion добавляет вопрос к блоку
func (b *Block) AddQuestion(question Question) {
	b.Questions = append(b.Questions, question)
}
