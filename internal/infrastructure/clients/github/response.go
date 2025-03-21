package github

import "time"

type RepoResponse struct {
	EventType string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *RepoResponse) GetDescription() string {
	return r.EventType
}

func (r *RepoResponse) GetDate() time.Time {
	return r.CreatedAt
}
