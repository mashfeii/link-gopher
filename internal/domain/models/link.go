package models

import (
	"net/url"
	"time"
)

type Link struct {
	LinkID     int64               `json:"link_id"`
	ChatID     int64               `json:"chat_id"`
	URL        string              `json:"url"`
	Tags       map[string]struct{} `json:"tags"`
	Filters    map[string][]string `json:"filters"`
	LastUpdate time.Time           `json:"last_update"`
}

func (l *Link) AddTag(tag string) {
	l.Tags[tag] = struct{}{}
}

func (l *Link) DeleteTag(tag string) {
	delete(l.Tags, tag)
}

func (l *Link) UpdateFilters(filters map[string][]string) {
	l.Filters = filters
}

func (l *Link) SetLastUpdate(lastUpdate time.Time) {
	l.LastUpdate = lastUpdate
}

func (l *Link) GetType() string {
	parsed, _ := url.Parse("https://" + l.URL)
	return parsed.Host
}
