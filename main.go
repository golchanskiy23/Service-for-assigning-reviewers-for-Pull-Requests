package main

import (
	//"Service-for-assigning-reviewers-for-Pull-Requests/internal/server"
	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	logger "Service-for-assigning-reviewers-for-Pull-Requests/pkg/logger"
	"os"
)

func main() {
	// config
	// logger
	l := logger.SetupLogger()
	if err := config.Configure(); err != nil {
		l.Error("failed to load .env file", "err", err)
		os.Exit(1)
	}
	l.Info("correctly loaded .env file")

	//srv := server.NewHTTPServer(":8080")
	//srv.Run()
	// logic
}
