package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/api/middleware"
	scrapper_api "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/api/servers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application/service"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/db"
	ormStorage "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage/orm"
	sqlStorage "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage/sql"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	configFileName := flag.String("config", "", "config file name/path")
	flag.Parse()

	cfg, err := config.NewConfig(*configFileName)
	if err != nil {
		slog.Error("unable to load config", slog.Any("error", err))
		os.Exit(1)
	}

	pool, err := db.NewPool(&cfg.Database)
	if err != nil {
		slog.Error("unable to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	var (
		linksRepo   repository.LinkRepository
		usersRepo   repository.UserRepository
		tagsRepo    repository.TagRepository
		filtersRepo repository.FilterRepository
	)

	if cfg.Database.AccessType == config.AccessTypeSQL {
		usersRepo = sqlStorage.NewSQLUserRepository(pool)
		linksRepo = sqlStorage.NewSQLLinkRepository(pool)
		tagsRepo = sqlStorage.NewSQLTagRepository(pool)
		filtersRepo = sqlStorage.NewSQLFilterRepository(pool)
	} else {
		usersRepo = ormStorage.NewORMUserRepository(pool)
		linksRepo = ormStorage.NewORMLinksRepository(pool)
		tagsRepo = ormStorage.NewORMTagsRepository(pool)
		filtersRepo = ormStorage.NewORMFiltersRepository(pool)
	}

	userService := service.NewUserService(usersRepo)
	linksService := service.NewLinkService(usersRepo, linksRepo, tagsRepo, filtersRepo)
	middlewares := []scrapper_api.MiddlewareFunc{
		middleware.SlogLogging,
	}
	server := servers.NewScrapperServer(cfg, userService, linksService, middlewares)

	schedulerDeps, err := application.NewDefaultDependencies(cfg, linksRepo)
	if err != nil {
		slog.Error("unable to create dependencies", slog.Any("error", err))
		os.Exit(1)
	}

	scheduler, err := application.StartScheduler(schedulerDeps)
	if err != nil {
		slog.Error("unable to load config", slog.Any("error", err))
		os.Exit(1)
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
