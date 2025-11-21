package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      App        `mapstructure:"app"`
	Server   HttpServer `mapstructure:"server"`
	Database DB         `mapstructure:"database"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"appversion"`
}

type DB struct {
	Name              string         `mapstructure:"name"`
	Port              int            `mapstructure:"port"`
	SSLMode           Mode           `mapstructure:"sslmode"`
	Schema            string         `mapstructure:"schema"`
	MaxPoolSize       int            `mapstructure:"maxpoolsize"`
	MaxConnLifetime   *time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnectTimeout *time.Duration `mapstructure:"max_connect_timeout"`
	QueryTimeout      *time.Duration `mapstructure:"query_timeout"`
}

type HttpServer struct {
	ReadTimeout     *time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    *time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration  `mapstructure:"shutdown_timeout"`
	Addr            string         `mapstructure:"addr"`
}

type Mode string

func NewConfig() (*Config, error) {
	cfg := &Config{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %s", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("marshaling error: %s", err)
	}
	return cfg, nil
}

