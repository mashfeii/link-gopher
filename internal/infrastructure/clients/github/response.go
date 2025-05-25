package github

import (
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
)

type RepoResponse struct {
	Title string `json:"title"`
	User  struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
}

func (p *RepoResponse) GetTitle() string {
	return p.Title
}

func (p *RepoResponse) GetUser() string {
	return p.User.Login
}

func (p *RepoResponse) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *RepoResponse) GetBody() string {
	return p.Body
}

type RepoResponseWrapper struct {
	Item RepoResponse
	Type models.EventType
}

func (w *RepoResponseWrapper) GetType() models.EventType { return w.Type }
func (w *RepoResponseWrapper) GetTitle() string          { return w.Item.GetTitle() }
func (w *RepoResponseWrapper) GetUser() string           { return w.Item.GetUser() }
func (w *RepoResponseWrapper) GetCreatedAt() time.Time   { return w.Item.GetCreatedAt() }
func (w *RepoResponseWrapper) GetBody() string           { return w.Item.GetBody() }
