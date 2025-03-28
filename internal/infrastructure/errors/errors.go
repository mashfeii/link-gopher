package errors

import "fmt"

type ErrUserNotFound struct {
	Code int
}

func NewErrUserNotFound() error {
	return ErrUserNotFound{Code: 404}
}

func (e ErrUserNotFound) Error() string {
	return "user not found"
}

type ErrLinkNotFound struct {
	Code int
}

func NewErrLinkNotFound() error {
	return ErrLinkNotFound{Code: 404}
}

func (e ErrLinkNotFound) Error() string {
	return "link not found"
}

type ErrUserAlreadyExists struct {
	UserID int64
}

func NewErrUserAlreadyExists(userID int64) error {
	return ErrUserAlreadyExists{
		UserID: userID,
	}
}

func (e ErrUserAlreadyExists) Error() string {
	return fmt.Sprintf("user with chatID %d already exists", e.UserID)
}

type ErrLinkAlreadyExists struct {
	URL string
}

func NewErrLinkAlreadyExists(url string) error {
	return ErrLinkAlreadyExists{
		URL: url,
	}
}

func (e ErrLinkAlreadyExists) Error() string {
	return fmt.Sprintf("link with URL %s already exists", e.URL)
}

type ErrInvalidFilterFormat struct {
	Filter string
}

func NewErrInvalidFilterFormat(filter string) error {
	return ErrInvalidFilterFormat{Filter: filter}
}

func (e ErrInvalidFilterFormat) Error() string {
	return fmt.Sprintf("filter must be in key:value format, got %s", e.Filter)
}

type ErrNoActiveSession struct {
	ChatID int64
}

func NewErrNoActiveSession(chatID int64) error {
	return ErrNoActiveSession{ChatID: chatID}
}

func (e ErrNoActiveSession) Error() string {
	return fmt.Sprintf("no active session for chatID %d", e.ChatID)
}

type ErrUnknownSessionType struct {
	ChatID int64
}

func NewErrUnknownSessionType(chatID int64) error {
	return ErrUnknownSessionType{ChatID: chatID}
}

func (e ErrUnknownSessionType) Error() string {
	return fmt.Sprintf("unknown session type for chatID %d", e.ChatID)
}

type ErrInvalidURL struct {
	URL string
}

func NewErrInvalidURL(url string) error {
	return ErrInvalidURL{URL: url}
}

func (e ErrInvalidURL) Error() string {
	return fmt.Sprintf("invalid URL: %s", e.URL)
}

type ErrTagNotFound struct {
	LinkID  int64
	TagName string
}

func NewErrTagNotFound(linkID int64, tagName string) error {
	return ErrTagNotFound{
		LinkID:  linkID,
		TagName: tagName,
	}
}

func (e ErrTagNotFound) Error() string {
	return fmt.Sprintf("tag %s not found for link %d", e.TagName, e.LinkID)
}

type ErrTagAlreadyExist struct {
	TagName string
	LinkID  int64
}

func NewErrTagAlreadyExist(tagName string, linkID int64) error {
	return ErrTagAlreadyExist{TagName: tagName, LinkID: linkID}
}

func (e ErrTagAlreadyExist) Error() string {
	return fmt.Sprintf("tag %s already exists for link %d", e.TagName, e.LinkID)
}

type ErrFilterAlreadyExists struct {
	LinkID int64
	Key    string
	Value  string
}

func NewErrFilterAlreadyExists(linkID int64, key, value string) error {
	return ErrFilterAlreadyExists{LinkID: linkID, Key: key, Value: value}
}

func (e ErrFilterAlreadyExists) Error() string {
	return fmt.Sprintf("filter with key %s and value %s already exists for link %d", e.Key, e.Value, e.LinkID)
}

type ErrFilterNotFound struct {
	LinkID int64
	Key    string
	Value  string
}

func NewErrFilterNotFound(linkID int64, key, value string) error {
	return ErrFilterNotFound{LinkID: linkID, Key: key, Value: value}
}

func (e ErrFilterNotFound) Error() string {
	return fmt.Sprintf("filter with key %s and value %s not found for link %d", e.Key, e.Value, e.LinkID)
}
