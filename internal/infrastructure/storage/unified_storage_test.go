package storage_test

import (
	"context"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage"
	"github.com/stretchr/testify/assert"
)

func TestCombinedRepository_CreateUser(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(123)

	t.Run("successful user creation", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)
	})

	t.Run("duplicate user error", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.ErrorContains(t, err, "already exists")
	})
}

func TestCombinedRepository_GetUser(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(124)

	t.Run("user not found error", func(t *testing.T) {
		_, err := repo.GetUser(ctx, userID)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("successful user retrieval", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		user, err := repo.GetUser(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, userID, user.ChatID)
	})
}

func TestCombinedRepository_DeleteUser(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(125)

	t.Run("user not found error", func(t *testing.T) {
		_, err := repo.DeleteUser(ctx, userID)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("successful user deletion", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		user, err := repo.DeleteUser(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, userID, user.ChatID)

		_, err = repo.GetUser(ctx, userID)
		assert.ErrorContains(t, err, "not found")
	})
}

func TestCombinedRepository_AddLink(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(126)
	testLink := &models.Link{
		URL:  "https://github.com/owner/repo",
		Tags: map[string]struct{}{"go": {}},
    Filters: map[string][]string{"filtername": {"filtervalue"}},
	}

	t.Run("successful link addition", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		err = repo.AddLink(ctx, userID, testLink)
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)

    user, err := repo.GetUser(ctx, userID)
    assert.NoError(t, err)

		assert.Len(t, links, 1)
		assert.Equal(t, testLink.URL, links[0].URL)
    assert.Equal(t, testLink.Tags, links[0].Tags)
    assert.Contains(t, user.Filters, "filtername")
    assert.Contains(t, user.Filters["filtername"], "filtervalue")
	})

	t.Run("duplicate link error", func(t *testing.T) {
		err := repo.AddLink(ctx, userID, testLink)
		assert.ErrorContains(t, err, "already exists")
	})

	t.Run("not found user error", func(t *testing.T) {
		err := repo.AddLink(ctx, 127, testLink)
		assert.ErrorContains(t, err, "not found")
	})
}

func TestCombinedRepository_DeleteLink(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(128)
	testLink := &models.Link{
		URL:  "https://github.com/owner/repo",
		Tags: map[string]struct{}{"go": {}},
	}

	t.Run("link not found error", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		_, err = repo.DeleteLink(ctx, userID, testLink.URL)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("user not found error", func(t *testing.T) {
		_, err := repo.DeleteLink(ctx, 129, testLink.URL)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("successful link deletion", func(t *testing.T) {
		err := repo.AddLink(ctx, userID, testLink)
		assert.NoError(t, err)

		_, err = repo.DeleteLink(ctx, userID, testLink.URL)
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, links, 0)
	})
}

func TestCombinedRepository_GetLinks(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(130)
	testLinks := []models.Link{
		{
			ChatID: userID,
			URL:    "github.com/owner/repo",
			Tags:   map[string]struct{}{"go": {}},
		},
		{
			ChatID: userID,
			URL:    "github.com/owner/repo2",
			Tags:   map[string]struct{}{"go": {}},
		},
	}

	t.Run("no links found", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, links, 0)
	})

	t.Run("user not found error", func(t *testing.T) {
		_, err := repo.GetLinks(ctx, 131)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("successful link retrieval", func(t *testing.T) {
		err := repo.AddLink(ctx, userID, &testLinks[0])
		assert.NoError(t, err)

		err = repo.AddLink(ctx, userID, &testLinks[1])
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, links, 2)
		assert.Contains(t, links, testLinks[0])
		assert.Contains(t, links, testLinks[1])
	})
}

func TestCombinedRepository_GetAllActiveLinks(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID1 := int64(132)
	userID2 := int64(133)
	testLinks := []models.Link{
		{
			ChatID: userID1,
			URL:    "github.com/owner/repo",
			Tags:   map[string]struct{}{"go": {}},
		},
		{
			ChatID: userID2,
			URL:    "github.com/owner/repo2",
			Tags:   map[string]struct{}{"go": {}},
		},
	}

	t.Run("no links found", func(t *testing.T) {
		links, err := repo.GetAllActiveLinks(ctx)
		assert.NoError(t, err)
		assert.Len(t, links, 0)
	})

	t.Run("successful link retrieval", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID1)
		assert.NoError(t, err)

		err = repo.CreateUser(ctx, userID2)
		assert.NoError(t, err)

		err = repo.AddLink(ctx, userID1, &testLinks[0])
		assert.NoError(t, err)

		err = repo.AddLink(ctx, userID2, &testLinks[1])
		assert.NoError(t, err)

		links, err := repo.GetAllActiveLinks(ctx)
		assert.NoError(t, err)

		assert.Len(t, links, 2)
		assert.Contains(t, links, testLinks[0])
		assert.Contains(t, links, testLinks[1])
	})
}

func TestCombinedRepository_AddTagToLink(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(134)
	testLink := &models.Link{
		URL:  "github.com/owner/repo",
		Tags: map[string]struct{}{"go": {}},
	}

	t.Run("link not found error", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		err = repo.AddTagToLink(ctx, userID, testLink.URL, "golang")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("user not found error", func(t *testing.T) {
		err := repo.AddTagToLink(ctx, 135, testLink.URL, "golang")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("successful tag addition", func(t *testing.T) {
		err := repo.AddLink(ctx, userID, testLink)
		assert.NoError(t, err)

		err = repo.AddTagToLink(ctx, userID, testLink.URL, "golang")
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, links, 1)
		assert.Contains(t, links[0].Tags, "golang")
	})
}

func TestCombinedRepository_DeleteTagFromLink(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(136)
	testLink := &models.Link{
		URL:  "github.com/owner/repo",
		Tags: map[string]struct{}{"go": {}, "golang": {}},
	}
	testLink2 := &models.Link{
		URL:  "github.com/owner/repo2",
		Tags: map[string]struct{}{"go": {}, "golang": {}},
	}

	t.Run("user not found error", func(t *testing.T) {
		err := repo.DeleteTagFromLink(ctx, 135, testLink.URL, "go")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("link not found error", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		err = repo.DeleteTagFromLink(ctx, userID, "github.com/owner/repo2", "go")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("tag not found error", func(t *testing.T) {
		err := repo.AddLink(ctx, userID, testLink)
		assert.NoError(t, err)

		err = repo.AddLink(ctx, userID, testLink2)
		assert.NoError(t, err)

		err = repo.DeleteTagFromLink(ctx, userID, testLink.URL, "python")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("successful tag deletion still used", func(t *testing.T) {
		err := repo.DeleteTagFromLink(ctx, userID, testLink.URL, "go")
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)

		user, err := repo.GetUser(ctx, userID)
		assert.NoError(t, err)

		assert.Len(t, links, 2)
		assert.NotContains(t, links[0].Tags, "go")
		assert.Contains(t, user.Tags, "go")
	})

	t.Run("successful tag deletion not used", func(t *testing.T) {
		err := repo.DeleteTagFromLink(ctx, userID, testLink2.URL, "go")
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)

		user, err := repo.GetUser(ctx, userID)
		assert.NoError(t, err)

		assert.Len(t, links, 2)
		assert.NotContains(t, links[1].Tags, "go")
		assert.NotContains(t, user.Tags, "go")
	})
}

func TestCombinedRepository_UpdateLinkFilters(t *testing.T) {
	repo := storage.NewCombinedRepository()
	ctx := context.Background()
	userID := int64(137)
	testLink := &models.Link{
		URL:     "github.com/owner/repo",
		Tags:    map[string]struct{}{"go": {}},
		Filters: map[string][]string{"filtername": {"filtervalue"}},
	}

	t.Run("link not found error", func(t *testing.T) {
		err := repo.CreateUser(ctx, userID)
		assert.NoError(t, err)

		err = repo.UpdateLinkFilters(ctx, userID, "github.com/owner/repo2", map[string][]string{"filtername": {"newfiltervalue"}})
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("user not found error", func(t *testing.T) {
		err := repo.UpdateLinkFilters(ctx, 138, testLink.URL, map[string][]string{"filtername": {"newfiltervalue"}})
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("successful filter update", func(t *testing.T) {
		err := repo.AddLink(ctx, userID, testLink)
		assert.NoError(t, err)

		filters := map[string][]string{"filtername": {"newfiltervalue"}, "filtername2": {"filtervalue2"}}
		err = repo.UpdateLinkFilters(ctx, userID, testLink.URL, filters)
		assert.NoError(t, err)

		links, err := repo.GetLinks(ctx, userID)
		assert.NoError(t, err)
		user, err := repo.GetUser(ctx, userID)
		assert.NoError(t, err)

		assert.Len(t, links, 1)
		assert.Equal(t, filters, links[0].Filters)

		assert.Contains(t, links[0].Filters["filtername"], "newfiltervalue")
		assert.Contains(t, user.Filters["filtername"], "newfiltervalue")
		assert.Contains(t, links[0].Filters["filtername2"], "filtervalue2")
	})
}
