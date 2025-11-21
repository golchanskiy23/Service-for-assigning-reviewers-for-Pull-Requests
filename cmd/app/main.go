package main

import (
	"fmt"
	"os"

	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/app"
	logger "Service-for-assigning-reviewers-for-Pull-Requests/pkg/logger"
)

const exitCodeError = 1

func main() {
	log := logger.SetupLogger()
	if err := config.SystemVarsInit(); err != nil {
		log.Error(fmt.Sprintf("failed to load .env file by error %v", err))
		os.Exit(exitCodeError)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error(fmt.Sprintf("failed to configuration by error %v", err))
		os.Exit(exitCodeError)
	}

	if err = app.Run(cfg, log); err != nil {
		log.Error(fmt.Sprintf("error running app: %v", err))
		os.Exit(exitCodeError)
	}
}
