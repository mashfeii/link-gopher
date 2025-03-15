package stackoverflow

import (
	"time"
)

type SOQuestionResponse struct {
	Items []struct {
		IsAnswered       bool  `json:"is_answered"`
		LastActivityDate int64 `json:"last_activity_date"`
	} `json:"items"`
}

func (r *SOQuestionResponse) GetDescription() string {
	if len(r.Items) == 0 {
		return ""
	}

	if r.Items[0].IsAnswered {
		return "Question is answered"
	}

	return "Still unanswered"
}

func (r *SOQuestionResponse) GetDate() time.Time {
	if len(r.Items) == 0 {
		return time.Time{}
	}

	return time.Unix(r.Items[0].LastActivityDate, 0)
}
