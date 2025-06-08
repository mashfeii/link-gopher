package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/api/servers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot"
)

func main() {
	configFileName := flag.String("config", "", "config file name")
	flag.Parse()

	cfg, err := config.NewConfig(*configFileName)
	if err != nil {
		slog.Error("unable to load config", slog.Any("error", err))
		os.Exit(1)
	}

	tgClient, err := telebot.NewBotClient(cfg)
	if err != nil {
		slog.Error("unable to create bot", slog.Any("error", err))
		os.Exit(1)
	}

	deps := application.NewBotDependencies(cfg, tgClient)
	server := servers.NewBotServer(deps)

	tgClient.Run()

	if err := server.ListenAndServe(); err != nil {
		slog.Error("unable to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
