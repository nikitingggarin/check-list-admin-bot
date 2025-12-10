package types

// CheckListData - интерфейс для данных чек-листа
type CheckListData interface {
	GetType() CheckListType
	GetID() int64
	GetName() string
	GetStatus() CheckListStatus
}

// SimpleCheckList - простой чек-лист
type SimpleCheckList struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Status    CheckListStatus `json:"status"`
	Questions []Question      `json:"questions"`
}

func (s *SimpleCheckList) GetType() CheckListType     { return TypeSimple }
func (s *SimpleCheckList) GetID() int64               { return s.ID }
func (s *SimpleCheckList) GetName() string            { return s.Name }
func (s *SimpleCheckList) GetStatus() CheckListStatus { return s.Status }

// NewDraftSimpleCheckList создает черновик простого чек-листа
func NewDraftSimpleCheckList(name string) *SimpleCheckList {
	return &SimpleCheckList{
		ID:        -1,
		Name:      name,
		Status:    StatusDraft,
		Questions: make([]Question, 0),
	}
}

// AddQuestion добавляет вопрос к простому чек-листу
func (c *SimpleCheckList) AddQuestion(question Question) {
	c.Questions = append(c.Questions, question)
}

// BlockedCheckList - чек-лист с блоками
type BlockedCheckList struct {
	ID     int64           `json:"id"`
	Name   string          `json:"name"`
	Status CheckListStatus `json:"status"`
	Blocks []Block         `json:"blocks"`
}

func (b *BlockedCheckList) GetType() CheckListType     { return TypeBlocked }
func (b *BlockedCheckList) GetID() int64               { return b.ID }
func (b *BlockedCheckList) GetName() string            { return b.Name }
func (b *BlockedCheckList) GetStatus() CheckListStatus { return b.Status }

// NewDraftBlockedCheckList создает черновик чек-листа с блоками
func NewDraftBlockedCheckList(name string) *BlockedCheckList {
	return &BlockedCheckList{
		ID:     -1,
		Name:   name,
		Status: StatusDraft,
		Blocks: make([]Block, 0),
	}
}

// AddBlock добавляет блок к чек-листу с блоками
func (c *BlockedCheckList) AddBlock(block Block) {
	c.Blocks = append(c.Blocks, block)
}
