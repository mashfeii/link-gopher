package telebot

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	scrapper_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

func (bot *BotClient) handleAuthorizationCheck(chatID int64) (*scrapper_client.GetLinksResponse, error) {
	resp, err := bot.scrapperClient.GetLinksWithResponse(context.TODO(),
		&scrapper_client.GetLinksParams{TgChatId: chatID},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to validate user: %w", err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "❌ Unauthorized, please use /start to register"))
		return nil, fmt.Errorf("unauthorized user")
	}

	return resp, nil
}

func (bot *BotClient) validateLink(chatID int64, url string) (bool, error) {
	resp, err := bot.scrapperClient.GetLinksWithResponse(context.TODO(),
		&scrapper_client.GetLinksParams{TgChatId: chatID},
	)
	if err != nil {
		return true, fmt.Errorf("unable to get links: %w", err)
	}

	if resp.JSON200 != nil && *resp.JSON200.Size > 0 {
		for _, link := range *resp.JSON200.Links {
			if *link.Url == url {
				return false, nil
			}
		}
	}

	_, _, gitErr := pkg.ValidateGithubURL(url)
	_, stackErr := pkg.ValidateStackOverflowURL(url)

	if gitErr != nil && stackErr != nil {
		return true, errors.NewErrInvalidURL(url)
	}

	return true, nil
}

func (bot *BotClient) extractFiltersAndTags(links []scrapper_client.LinkResponse) (filters, tags []string) {
	availableFilters := []string{}
	availableTags := make(map[string]struct{})

	for _, link := range links {
		for _, filter := range *link.Filters {
			if !slices.Contains(availableFilters, filter) {
				availableFilters = append(availableFilters, filter)
			}
		}

		for _, tag := range *link.Tags {
			availableTags[tag] = struct{}{}
		}
	}

	return availableFilters, lo.Keys(availableTags)
}

func (bot *BotClient) filterLinks(
	links []scrapper_client.LinkResponse,
	hasTags bool,
	selectedTags []string,
	filterValid bool,
	filter string,
) []scrapper_client.LinkResponse {
	if !hasTags && !filterValid {
		return links
	}

	filtered := make([]scrapper_client.LinkResponse, 0)

	for _, link := range links {
		matchTags := true
		if hasTags {
			matchTags = hasAllTags(link, selectedTags)
		}

		matchFilter := false
		if filterValid && matchTags {
			matchFilter = slices.Contains(*link.Filters, filter)
		}

		if matchTags && matchFilter {
			filtered = append(filtered, link)
		}
	}

	return filtered
}

func validateLinksResponse(chatID int64, bot *BotClient, resp *scrapper_client.GetLinksResponse) (*[]scrapper_client.LinkResponse, error) {
	if resp.JSON200 == nil || *resp.JSON200.Size == 0 {
		_, err := bot.bot.Send(tgbotapi.NewMessage(chatID, "😔 No links found, use /track to add a new link"))
		return nil, err
	}

	return resp.JSON200.Links, nil
}

func hasAllTags(link scrapper_client.LinkResponse, tags []string) bool {
	for _, tag := range tags {
		if !slices.Contains(*link.Tags, tag) {
			return false
		}
	}

	return true
}

func handleErrorResponse(response *scrapper_client.GetLinksResponse) error {
	if response.JSON400 != nil {
		return fmt.Errorf("invalid request: %s", *response.JSON400.Description)
	}

	return nil
}
