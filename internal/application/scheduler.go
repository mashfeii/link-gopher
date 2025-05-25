package application

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-co-op/gocron/v2"

	bot_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/bot"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
)

func StartScheduler(deps *SchedulerDependencies) (gocron.Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(time.Duration(deps.Config.Serving.Interval)*time.Minute),
		gocron.NewTask(
			checkUpdates,
			deps.Repo,
			deps.BotClient,
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
	repo repository.LinkRepository,
	botClient bot_client.ClientInterface,
	ghClient models.LinkChecker,
	soClient models.LinkChecker,
) {
	const op = "application.checkUpdates"

	links, err := repo.GetAllActiveLinks(context.Background())
	if err != nil {
		return
	}

	for _, link := range links {
		var client models.LinkChecker

		switch link.GetType() {
		case models.LinkTypeGithub:
			client = ghClient
		case models.LinkTypeStackOverflow:
			client = soClient
		case models.LinkTypeUnknown:
			slog.Warn("unknown link type, skipping", "operation", op, "link", link.URL)
			continue
		}

		updates, err := client.GetUpdates(link.URL, link.LastUpdate)
		if err != nil {
			slog.Error("failed to get updates", "operation", op, "link", link.URL, "error", err)
			continue
		}

		if len(updates) == 0 {
			slog.Info("no updates found", "operation", op, "link", link.URL)
			continue
		}

		for _, update := range updates {
			if update.GetCreatedAt().Before(link.LastUpdate) {
				continue
			}

			var message strings.Builder

			if link.GetType() == models.LinkTypeStackOverflow {
				title, _ := client.GetQuestionTitle(link.URL)
				message.WriteString(fmt.Sprintf("New %s on question %q: %s by %s at %s",
					update.GetType(), title, update.GetTitle(), update.GetUser(), update.GetCreatedAt()))
			} else {
				message.WriteString(fmt.Sprintf("New %s: %s by %s at %s",
					update.GetType(), update.GetTitle(), update.GetUser(), update.GetCreatedAt()))
			}

			finalResponse := message.String()

			resp, err := botClient.PostUpdates(context.TODO(), bot_client.LinkUpdate{
				TgChatId:    &link.ChatID,
				Url:         &link.URL,
				Description: &finalResponse,
			})
			if err != nil {
				slog.Error("failed to post update", "opeartion", op, "link", link.URL, "error", err)
				continue
			}

			if resp != nil {
				resp.Body.Close()
			}
		}

		// TODO: expand repository interface to support update operation.
		link.SetLastUpdate(time.Now())

		if len(updates) > 0 {
			link.SetLastUpdate(updates[len(updates)-1].GetCreatedAt())
		}
	}
}
