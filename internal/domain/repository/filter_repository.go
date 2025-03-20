package repository

import "context"

type FilterRepository interface {
	AddFilterToLink(ctx context.Context, linkID int64, key, value string) error
	RemoveFiltersFromLink(ctx context.Context, linkID int64) error
	GetFiltersByLink(ctx context.Context, linkID int64) ([]string, error)
}
