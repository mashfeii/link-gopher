package service

import (
	"context"
	"fmt"
	"net/http"

	scrapperclient "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type ScrapperClientWrapper struct {
	client scrapperclient.ClientWithResponsesInterface
}

func NewScrapperClientWrapper(client scrapperclient.ClientWithResponsesInterface) *ScrapperClientWrapper {
	return &ScrapperClientWrapper{
		client: client,
	}
}

func (s *ScrapperClientWrapper) RegisterUser(ctx context.Context) error {
	const op = "scrapperClientWrapper.RegisterUser"
	chatID := ctx.Value(core.ContextKeyChatID).(int64)

	resp, err := s.client.PostTgChatIdWithResponse(ctx, chatID)
	if err != nil {
		return fmt.Errorf("%s: unable to register user: %w", op, err)
	}

	if resp.StatusCode() != http.StatusOK || resp.StatusCode() == http.StatusAlreadyReported {
		return errors.NewErrUserAlreadyExists(chatID)
	}

	return nil
}

func (s *ScrapperClientWrapper) GetUser(ctx context.Context) (bool, error) {
	const op = "scrapperClientWrapper.GetUser"
	chatID := ctx.Value(core.ContextKeyChatID).(int64)

	resp, err := s.client.GetLinksWithResponse(ctx, &scrapperclient.GetLinksParams{TgChatId: chatID})
	if err != nil {
		return false, fmt.Errorf("%s: unable to get user: %w", op, err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return false, errors.NewErrUserNotFound(chatID)
	}

	return true, nil
}

func (s *ScrapperClientWrapper) SaveLink(ctx context.Context, linkData models.LinkData) error {
	const op = "scrapperClientWrapper.SaveLink"
	chatID := ctx.Value(core.ContextKeyChatID).(int64)

	params := &scrapperclient.PostLinksParams{
		TgChatId: chatID,
	}
	body := scrapperclient.PostLinksJSONRequestBody{
		Link:    &linkData.URL,
		Tags:    &linkData.Tags,
		Filters: &linkData.Filters,
	}

	resp, err := s.client.PostLinksWithResponse(ctx, params, body)
	if err != nil {
		return fmt.Errorf("%s: unable to save link: %w", op, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%s: unable to save link: %s", op, resp.Status())
	}

	return nil
}

func (s *ScrapperClientWrapper) GetLinks(ctx context.Context) ([]models.LinkData, error) {
	const op = "scrapperClientWrapper.GetLinks"
	chatID := ctx.Value(core.ContextKeyChatID).(int64)

	params := &scrapperclient.GetLinksParams{
		TgChatId: chatID,
	}

	resp, err := s.client.GetLinksWithResponse(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to get links: %w", op, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("%s: unable to get links: %s", op, resp.Status())
	}

	var links []models.LinkData
	for _, link := range *resp.JSON200.Links {
		links = append(links, models.LinkData{
			URL:     *link.Url,
			Tags:    *link.Tags,
			Filters: *link.Filters,
		})
	}

	return links, nil
}

func (s *ScrapperClientWrapper) DeleteLink(ctx context.Context) error {
	const op = "scrapperClientWrapper.DeleteLink"
	chatID := ctx.Value(core.ContextKeyChatID).(int64)
	url := ctx.Value(core.ContextLinkURL).(string)

	params := &scrapperclient.DeleteLinksParams{
		TgChatId: chatID,
	}
	body := scrapperclient.DeleteLinksJSONRequestBody{
		Link: &url,
	}

	resp, err := s.client.DeleteLinksWithResponse(ctx, params, body)
	if err != nil {
		return fmt.Errorf("%s: unable to delete link: %w", op, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%s: unable to delete link: %s", op, resp.Status())
	}

	return nil
}
