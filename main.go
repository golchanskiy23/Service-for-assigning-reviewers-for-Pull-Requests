package main

import (
	//"Service-for-assigning-reviewers-for-Pull-Requests/internal/server"
	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	"Service-for-assigning-reviewers-for-Pull-Requests/internal/app"
	logger "Service-for-assigning-reviewers-for-Pull-Requests/pkg/logger"
	"os"
)

func main() {
	l := logger.SetupLogger()
	if err := config.SystemVarsInit(); err != nil {
		l.Error("failed to load .env file", "err", err)
		os.Exit(1)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		l.Error("failed to configuration", "err", err)
		os.Exit(1)
	}

	if err = app.Run(cfg); err != nil {
		l.Error("Error running app: %v", err)
		os.Exit(1)
	}

	//srv := server.NewHTTPServer(":8080")
	//srv.Run()
	// logic
}
