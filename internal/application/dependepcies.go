package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	bot_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/bot"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/clients/github"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/clients/stackoverflow"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot"
)

type (
	SchedulerDependencies struct {
		Repo                repository.LinkRepository
		BotClient           bot_client.ClientInterface
		GithubClient        models.LinkChecker
		StackoverflowClient models.LinkChecker
		Config              *config.Config
	}
	BotDependencies struct {
		TgClient *telebot.BotClient
		Config   *config.Config
	}
)

func NewDefaultDependencies(cfg *config.Config, repo repository.LinkRepository) (*SchedulerDependencies, error) {
	client, err := bot_client.NewClient(fmt.Sprintf("http://%s:%d", cfg.Serving.Host, cfg.Serving.BotPort))
	if err != nil {
		return nil, err
	}

	defaultClient := http.Client{
		Timeout: time.Duration(cfg.Serving.Interval) * time.Second,
	}
	githubClient := github.NewClient(
		cfg.Secret.GitHubToken,
		"https://api.github.com/repos",
		&defaultClient,
	)
	stackoverflowClient := stackoverflow.NewClient(
		cfg.Secret.StackOverflowToken,
		"https://api.stackexchange.com/2.3",
		&defaultClient,
	)

	return &SchedulerDependencies{
		Repo:                repo,
		BotClient:           client,
		GithubClient:        githubClient,
		StackoverflowClient: stackoverflowClient,
		Config:              cfg,
	}, nil
}

func NewBotDependencies(cfg *config.Config, client *telebot.BotClient) *BotDependencies {
	return &BotDependencies{
		Config:   cfg,
		TgClient: client,
	}
}
