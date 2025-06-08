package models

import domainModels "github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"

type LinkResponse struct {
	LinkID      int64  `json:"link_id"`
	ChatID      int64  `json:"chat_id"`
	URL         string `json:"url"`
	LastUpdated string `json:"last_update"`
}

func ToLinkResponse(link *domainModels.Link) *LinkResponse {
	return &LinkResponse{
		LinkID:      link.LinkID,
		ChatID:      link.ChatID,
		URL:         link.URL,
		LastUpdated: link.LastUpdated.Format("2006-01-02 15:04:05"),
	}
}
