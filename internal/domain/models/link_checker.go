package models

import "time"

type EventType string

const (
	EventTypeAnswer      EventType = "Answer"
	EventTypeComment     EventType = "Comment"
	EventTypePullRequest EventType = "PullRequest"
	EventTypeIssue       EventType = "Issue"
)

type Event interface {
	GetType() EventType
	GetTitle() string
	GetUser() string
	GetCreatedAt() time.Time
	GetBody() string
}

type LinkChecker interface {
	GetUpdates(url string, since time.Time) ([]Event, error)
	GetQuestionTitle(url string) (string, error)
}
