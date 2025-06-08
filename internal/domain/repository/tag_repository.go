package repository

import (
	"context"
)

type TagRepository interface {
	AddTagToLink(ctx context.Context, chatID int64, linkURL, tagName string) error
	AddTagToUser(ctx context.Context, userID int64, tagName string) (*int64, error)
	GetTagsByLink(ctx context.Context, linkID int64) ([]string, error)
	GetTagsByUser(ctx context.Context, chatID int64) ([]string, error)
	RemoveTagsFromLink(ctx context.Context, linkID int64) error
}
