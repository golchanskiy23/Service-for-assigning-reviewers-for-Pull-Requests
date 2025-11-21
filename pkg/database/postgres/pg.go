package postgres

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	"Service-for-assigning-reviewers-for-Pull-Requests/util"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

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
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		"postgres",
		os.Getenv("POSTGRES_UNSAFE_USERNAME"),
		os.Getenv("POSTGRES_UNSAFE_PASSWORD"),
		os.Getenv("DB_HOST"),
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
}

func (s *DatabaseSource) Close() {
	s.Pool.Close()
}

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
	cfg.MaxConns = int32(src.MaxPoolSize)
	ctx := context.Background()
	for attempt := 0; attempt < src.MaxConnectAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			src.Pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
			if err == nil {
				fmt.Println(src)
				return src, nil
			}
		}
		if attempt == src.MaxConnectAttempts {
			return nil, errors.New("max connection attempts exceeded; connection is failed!")
		}

		jitter := util.CreateNewDelay(attempt, src.MaxConnectTimeout)
		time.Sleep(jitter)
	}
	return nil, ctx.Err()
}
