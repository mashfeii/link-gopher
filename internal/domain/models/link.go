package models

import (
	"net/url"
	"time"
)

type LinkType string

const (
	LinkTypeGithub        LinkType = "github.com"
	LinkTypeStackOverflow LinkType = "stackoverflow.com"
	LinkTypeUnknown       LinkType = "unknown"
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

func (l *Link) GetType() LinkType {
	parsed, err := url.Parse("https://" + l.URL)
	if err != nil {
		return LinkTypeUnknown
	}

	return LinkType(parsed.Host)
}
