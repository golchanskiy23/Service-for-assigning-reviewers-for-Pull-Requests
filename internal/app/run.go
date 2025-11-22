package app

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/handlers"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/server"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"time"
)

func initPostgres(cfg *config.Config) (*database.DatabaseSource, error) {
	opts := []database.Option{
		database.SetMaxPoolSize(cfg.Database.MaxPoolSize),
	}

	if cfg.Database.MaxConnLifetime != nil {
		opts = append(opts, database.SetMaxConnLifetime(*cfg.Database.MaxConnLifetime))
	}
	if cfg.Database.MaxConnectTimeout != nil {
		opts = append(opts, database.SetMaxConnectTimeout(*cfg.Database.MaxConnectTimeout))
	}

	db, err := database.NewStorage(
		database.GetConnection(&cfg.Database),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres: %w", err)
	}
	if err = db.Pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}
	return db, nil
}

func initDBRepository(db *database.DatabaseSource) *postgres.Repository {
	return postgres.CreateNewDBRepository(db)
}

func Run(cfg *config.Config, logger *slog.Logger) error {
	db, err := initPostgres(cfg)
	if err != nil {
		return err
	}
	defer func(db *database.DatabaseSource) {
		db.Close()
	}(db)

	pgRepository := initDBRepository(db)
	s := handlers.CreateNewService(pgRepository)
	r := chi.NewMux()
	server.RegisterRoutes(s, r)
	srv := server.StartServer(cfg, r, logger)

	go func() {
		<-time.After(cfg.Server.ShutdownTimeout)
		logger.Info("shutdown timeout reached")
		srv.FullShutdownTimeout(logger)
	}()

	srv.GracefulShutdown(logger)
	return nil
}
