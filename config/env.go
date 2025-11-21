package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func SystemVarsInit() error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}
