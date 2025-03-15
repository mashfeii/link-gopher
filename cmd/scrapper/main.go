package main

import (
	"flag"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/api/servers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
)

func main() {
	configFileName := flag.String("config", "", "config file name/path")
	flag.Parse()

	cfg, err := config.NewConfig(*configFileName)
	if err != nil {
		slog.Error("unable to load config", slog.Any("error", err))
		return
	}

	deps, err := application.NewScrapperDependencies(cfg)
	if err != nil {
		slog.Error("unable to create dependencies", slog.Any("error", err))
		return
	}

	server := servers.NewScrapperServer(deps)

	scheduler, err := application.StartScheduler(deps)
	if err != nil {
		slog.Error("unable to load config", slog.Any("error", err))
		return
	}

	defer func() {
		err := scheduler.Shutdown()
		if err != nil {
			slog.Error("unable to shutdown scheduler", slog.Any("error", err))
		}
	}()

	scheduler.Start()

	if err := server.ListenAndServe(); err != nil {
		slog.Error("unable to start server", slog.Any("error", err))
		return
	}
}
