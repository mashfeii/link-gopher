package models

import "time"

type UpdateType string

const (
	UpdateTypeAnswer      UpdateType = "answer"
	UpdateTypeComment     UpdateType = "comment"
	UpdateTypeIssue       UpdateType = "issue"
	UpdateTypePullRequest UpdateType = "pull_request"
)

type Update struct {
	ID          int64
	LinkID      int64
	UpdateType  UpdateType
	Title       string
	UserName    string
	BodyPreview string
	CreatedAt   time.Time
	Sent        bool
}
