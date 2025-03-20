package models

import (
	"net/url"
	"time"
)

type Link struct {
	LinkID     int64     `json:"link_id"`
	ChatID     int64     `json:"chat_id"`
	URL        string    `json:"url"`
	LastUpdate time.Time `json:"last_update"`
}

func (l *Link) SetLastUpdate(update time.Time) {
	l.LastUpdate = update
}

func (l *Link) GetType() string {
	parsed, _ := url.Parse("https://" + l.URL)
	return parsed.Host
}
