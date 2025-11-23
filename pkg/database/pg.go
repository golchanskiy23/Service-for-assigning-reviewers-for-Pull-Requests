package database

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/pkg/util"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"Service-for-assigning-reviewers-for-Pull-Requests/config"
)

//nolint:revive // exported: DatabaseSource is a clear and descriptive name
type DatabaseSource struct {
	Pool               *pgxpool.Pool
	MaxPoolSize        int
	MaxConnectTimeout  time.Duration
	MaxConnLifetime    time.Duration
	MaxConnectAttempts int
}

const (
	defaultMaxPoolSize       = 5
	defaultMaxConnLifetime   = 600 * time.Second
	defaultMaxConnectTimeout = 1 * time.Second
	defaultMaxConnAttempts   = 5
)

func GetConnection(cfg *config.DB) string {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = ""
	}
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		"postgres",
		os.Getenv("POSTGRES_UNSAFE_USERNAME"),
		os.Getenv("POSTGRES_UNSAFE_PASSWORD"),
		host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	return dsn
}

func (s *DatabaseSource) Close() {
	s.Pool.Close()
}

//nolint:revive // Complex initialization logic
func NewStorage(url string, options ...Option) (*DatabaseSource, error) {
	src := &DatabaseSource{
		MaxPoolSize:        defaultMaxPoolSize,
		MaxConnLifetime:    defaultMaxConnLifetime,
		MaxConnectTimeout:  defaultMaxConnectTimeout,
		MaxConnectAttempts: defaultMaxConnAttempts,
	}
	for _, option := range options {
		option(src)
	}

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("creation db storage error: %w", err)
	}

	//nolint:gosec // int to int32 conversion is safe for MaxPoolSize values
	cfg.MaxConns = int32(src.MaxPoolSize)
	cfg.MaxConnLifetime = src.MaxConnLifetime
	const (
		healthCheckPeriod = 30 * time.Second
		maxConnIdleTime   = 5 * time.Minute
	)
	cfg.HealthCheckPeriod = healthCheckPeriod
	cfg.MaxConnIdleTime = maxConnIdleTime

	ctx := context.Background()

	for attempt := range src.MaxConnectAttempts {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			src.Pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
			if err == nil {
				return src, nil
			}
		}

		if attempt == src.MaxConnectAttempts {
			return nil, errors.New("max connection attempts exceeded")
		}

		jitter := util.CreateNewDelay(attempt, src.MaxConnectTimeout)
		time.Sleep(jitter)
	}

	return nil, ctx.Err()
}
