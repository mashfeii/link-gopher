package application

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	bot_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/bot"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/clients/github"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/clients/stackoverflow"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot"
)

type (
	ScrapperDependencies struct {
		Repo                repository.UnifiedRepository
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

func NewScrapperDependencies(cfg *config.Config) (*ScrapperDependencies, error) {
	botServer := fmt.Sprintf("http://%s:%d", cfg.Serving.Host, cfg.Serving.BotPort)

	botClient, err := bot_client.NewClient(botServer)
	if err != nil {
		return nil, err
	}

	return &ScrapperDependencies{
		Repo:                storage.NewCombinedRepository(),
		BotClient:           botClient,
		GithubClient:        github.NewClient(cfg.Secret.GitHubToken),
		StackoverflowClient: stackoverflow.NewClient(cfg.Secret.StackOverflowToken),
		Config:              cfg,
	}, nil
}

func NewBotDependencies(cfg *config.Config, client *telebot.BotClient) *BotDependencies {
	return &BotDependencies{
		Config:   cfg,
		TgClient: client,
	}
}
