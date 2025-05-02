package models

type Session interface {
	GetState() State
	GetLastBotMessageID() *int
	GetLastUserMessageID() *int
}

type BaseSession struct {
	State             State
	LastBotMessageID  *int // Message to edit
	LastUserMessageID *int // Message to remove
}

func (s *BaseSession) GetState() State {
	return s.State
}

func (s *BaseSession) GetLastBotMessageID() *int {
	return s.LastBotMessageID
}

func (s *BaseSession) GetLastUserMessageID() *int {
	return s.LastUserMessageID
}

type TrackSession struct {
	BaseSession

	URL     string
	Tags    []string
	Filters []string
}

type ListUntrackSession struct {
	BaseSession

	AvailableTags []string
	SelectedTags  []string

	AvailableFilters []string
	FilterName       string
	FilterValue      string
}
