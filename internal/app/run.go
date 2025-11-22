package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/handlers"

	"github.com/go-chi/chi/v5"

	//nolint:revive // dependency
	postgres "Service-for-assigning-reviewers-for-Pull-Requests/internal/repository/postgres"
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/database"
	server "Service-for-assigning-reviewers-for-Pull-Requests/pkg/server"
)

const (
	zero = 0
)

func initPostgres(cfg *config.Config) (*database.DatabaseSource, error) {
	opts := []database.Option{
		database.SetMaxPoolSize(cfg.Database.MaxPoolSize),
	}

	if cfg.Database.MaxConnLifetime != nil {
		opts = append(
			opts,
			database.SetMaxConnLifetime(*cfg.Database.MaxConnLifetime),
		)
	}

	if cfg.Database.MaxConnectTimeout != nil {
		opts = append(
			opts,
			database.SetMaxConnectTimeout(*cfg.Database.MaxConnectTimeout),
		)
	}

	db, err := database.NewStorage(
		database.GetConnection(&cfg.Database),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create internal storage: %w", err)
	}

	ctx := context.Background()
	var pingErr error
	const maxPingAttempts = 10
	const retryDelay = 300 * time.Millisecond
	for attempt := zero; attempt < maxPingAttempts; attempt++ {
		pingErr = db.Pool.Ping(ctx)
		if pingErr == nil {
			break
		}
		wait := time.Duration(attempt+1) * retryDelay
		time.Sleep(wait)
	}
	if pingErr != nil {
		return nil, fmt.Errorf("ping error: %w", pingErr)
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
	s := handlers.CreateNewService(pgRepository, logger)
	r := chi.NewMux()
	server.RegisterRoutes(s, r)
	srv := server.StartServer(cfg, r, logger)

	go func() {
		<-time.After(cfg.Server.ShutdownTimeout)
		logger.Info("shutdown timeout reached")
		if err := srv.FullShutdownTimeout(logger); err != nil {
			logger.Error("failed to shutdown server", "error", err)
		}
	}()

	srv.GracefulShutdown(logger)

	return nil
}
