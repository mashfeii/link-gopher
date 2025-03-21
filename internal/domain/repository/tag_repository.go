package repository

import (
	"context"
)

type TagRepository interface {
	AddTagToLink(ctx context.Context, linkID int64, tagName string) error
	RemoveTagsFromLink(ctx context.Context, linkID int64) error
	GetTagsByLink(ctx context.Context, linkID int64) ([]string, error)
}
