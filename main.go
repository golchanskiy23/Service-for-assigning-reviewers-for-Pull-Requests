package main

import (
	//"Service-for-assigning-reviewers-for-Pull-Requests/internal/server"
	logger "Service-for-assigning-reviewers-for-Pull-Requests/pkg/logger"
)

func main() {
	// config
	// logger
	l := logger.SetupLogger()
	l.Info("starting service")
	//srv := server.NewHTTPServer(":8080")
	//srv.Run()
	// logic
}
