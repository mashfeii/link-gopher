package stackoverflow

import (
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
)

type SQLQuestionUser struct {
	DisplayName string `json:"display_name"`
}

type SQLQuestionItem struct {
	Title        string          `json:"title"`
	User         SQLQuestionUser `json:"owner"`
	CreationDate int64           `json:"creation_date"`
	Body         string          `json:"body"`
}

type SQLQuestionResponse struct {
	Items []SQLQuestionItem `json:"items"`
}

func (q *SQLQuestionItem) GetTitle() string {
	return q.Title // Assuming the body contains the title, adjust as necessary
}

func (q *SQLQuestionItem) GetUser() string {
	return q.User.DisplayName
}

func (q *SQLQuestionItem) GetCreatedAt() time.Time {
	return time.Unix(q.CreationDate, 0)
}

func (q *SQLQuestionItem) GetBody() string {
	return q.Body[:200]
}

type SQLQuestionItemWrapper struct {
	Item SQLQuestionItem
	Type models.EventType
}

func (w *SQLQuestionItemWrapper) GetType() models.EventType { return w.Type }
func (w *SQLQuestionItemWrapper) GetTitle() string          { return w.Item.GetTitle() }
func (w *SQLQuestionItemWrapper) GetUser() string           { return w.Item.GetUser() }
func (w *SQLQuestionItemWrapper) GetCreatedAt() time.Time   { return w.Item.GetCreatedAt() }
func (w *SQLQuestionItemWrapper) GetBody() string           { return w.Item.GetBody() }
