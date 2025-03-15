package models

import (
	"math/rand/v2"
	"slices"
)

type User struct {
	ChatID  int64               `json:"chat_id"`
	Links   map[string]Link     `json:"user_links"`
	Tags    map[string]struct{} `json:"user_tags"`
	Filters map[string][]string `json:"user_filters"`
}

func (u *User) AddLink(link *Link) {
	link.LinkID = rand.Int64() //nolint:gosec // repositroy step
	u.Links[link.URL] = *link

	for tag := range link.Tags {
		u.Tags[tag] = struct{}{}
	}

	for key, value := range link.Filters {
		if _, ok := u.Filters[key]; !ok {
			u.Filters[key] = value
			continue
		}

		for _, filter := range value {
			if !slices.Contains(u.Filters[key], filter) {
				u.Filters[key] = append(u.Filters[key], filter)
			}
		}
	}
}

func (u *User) DeleteLink(link *Link) {
	delete(u.Links, link.URL)

	for tag := range link.Tags {
		stillUsed := false

		for _, l := range u.Links {
			if _, ok := l.Tags[tag]; ok {
				stillUsed = true
				break
			}
		}

		if !stillUsed {
			delete(u.Tags, tag)
		}
	}

	for key, filters := range link.Filters {
		for _, filter := range filters {
			for i, f := range u.Filters[key] {
				if f == filter {
					u.Filters[key] = append(u.Filters[key][:i], u.Filters[key][i+1:]...)
					break
				}
			}
		}
	}
}

func (u *User) AddTag(tag string) {
	u.Tags[tag] = struct{}{}
}

func (u *User) DeleteTag(tag string) {
	delete(u.Tags, tag)
}
