package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	bot_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/bot"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
	"github.com/go-co-op/gocron/v2"
)

func StartScheduler(deps *ScrapperDependencies) (gocron.Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(deps.Config.Serving.Interval*time.Minute),
		gocron.NewTask(
			checkUpdates,
			deps.BotClient,
			deps.Repo,
			deps.GithubClient,
			deps.StackoverflowClient,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler job: %w", err)
	}

	scheduler.Start()

	return scheduler, nil
}

func checkUpdates(
	linksRepo repository.UnifiedRepository,
	botClient bot_client.ClientInterface,
	ghClient models.LinkChecker,
	soClient models.LinkChecker,
) {
	links, err := linksRepo.GetAllActiveLinks(context.Background())
	if err != nil {
		return
	}

	for _, link := range links {
		switch link.GetType() {
		case "github.com":
			handleEvent(&link, ghClient, botClient)

		case "stackoverflow.com":
			handleEvent(&link, soClient, botClient)
		}
	}
}

func handleEvent(link *models.Link, client models.LinkChecker, bot bot_client.ClientInterface) {
	event, err := client.GetEvent(link.URL)
	if err != nil {
		slog.Error("unable to get repo events", slog.Any("error", err))
		return
	}

	if event.GetDate().Before(link.LastUpdate) {
		return
	}

	var (
		chatID      = link.ChatID
		url         = link.URL
		description = event.GetDescription()
	)

	resp, err := bot.PostUpdates(context.TODO(), bot_client.LinkUpdate{
		TgChatId:    &chatID,
		Description: &description,
		Url:         &url,
	})
	if err != nil {
		slog.Error("unable to send update", slog.Any("error", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unable to send update", slog.Any("status", resp.StatusCode))
		return
	}

	link.SetLastUpdate(event.GetDate())
}
