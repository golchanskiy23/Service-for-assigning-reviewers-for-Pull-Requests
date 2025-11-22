package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   HTTPServer `mapstructure:"server"`
	App      App        `mapstructure:"app"`
	Database DB         `mapstructure:"database"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"appversion"`
}

type DB struct {
	MaxConnLifetime   *time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnectTimeout *time.Duration `mapstructure:"max_connect_timeout"`
	QueryTimeout      *time.Duration `mapstructure:"query_timeout"`
	Name              string         `mapstructure:"name"`
	SSLMode           Mode           `mapstructure:"sslmode"`
	Schema            string         `mapstructure:"schema"`
	Port              int            `mapstructure:"port"`
	MaxPoolSize       int            `mapstructure:"maxpoolsize"`
}

type HTTPServer struct {
	ReadTimeout     *time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    *time.Duration `mapstructure:"write_timeout"`
	Addr            string         `mapstructure:"addr"`
	ShutdownTimeout time.Duration  `mapstructure:"shutdown_timeout"`
}

type Mode string

func NewConfig() (*Config, error) {
	cfg := &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("marshaling error: %w", err)
	}

	return cfg, nil
}
