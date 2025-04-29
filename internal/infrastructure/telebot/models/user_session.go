package models

import "slices"

type (
	State       int
	SessionType int
)

const (
	SessionTypeTrack = iota
	SessionTypeListUntrack
)

const (
	StateWaitingURL = iota
	StateWaitingTags
	StateWaitingFilters
	StateWaitingTagsSelection
	StateWaitingFiltersSelection
)

type UserSession struct {
	Type    SessionType
	State   State
	Command string

	URL     string
	Tags    []string
	Filters []string

	SelectedTags       []string
	AvailableTags      []string
	AvailableFilters   []string
	CurrentFilterName  string
	CurrentFilterValue string
	LastMessageID      *int
	LastUserMessageID  *int
}

func (s *UserSession) ToggleTag(tag string) {
	for i, selectedTag := range s.SelectedTags {
		if selectedTag == tag {
			s.SelectedTags = slices.Delete(s.SelectedTags, i, i+1)
			return
		}
	}

	s.SelectedTags = append(s.SelectedTags, tag)
}
