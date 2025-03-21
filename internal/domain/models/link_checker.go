package models

import "time"

type Event interface {
	GetDescription() string
	GetDate() time.Time
}

type LinkChecker interface {
	GetEvent(url string) (Event, error)
}
