package main

import (
	//"Service-for-assigning-reviewers-for-Pull-Requests/internal/server"
	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/app"
	logger "Service-for-assigning-reviewers-for-Pull-Requests/pkg/logger"
	"fmt"
	"os"
)

func main() {
	l := logger.SetupLogger()
	if err := config.SystemVarsInit(); err != nil {
		l.Error(fmt.Sprintf("failed to load .env file by error %v", err))
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		l.Error(fmt.Sprintf("failed to configuration by error %v", err))
		os.Exit(1)
	}

	if err = app.Run(cfg, l); err != nil {
		l.Error(fmt.Sprintf("error running app: %v", err))
		os.Exit(1)
	}
}
