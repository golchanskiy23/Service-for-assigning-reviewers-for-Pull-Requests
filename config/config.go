package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func Configure() error {
	if err := godotenv.Load(".env.example"); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}
