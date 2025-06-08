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
	LinkID      int64
	ChatID      int64
	URL         string
	Type        LinkType
	LastChecked time.Time
	LastUpdated time.Time
}

func (l *Link) GetType() LinkType {
	parsed, err := url.Parse("https://" + l.URL)
	if err != nil {
		return LinkTypeUnknown
	}

	return LinkType(parsed.Host)
}
