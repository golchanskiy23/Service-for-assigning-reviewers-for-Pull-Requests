package main

import "Service-for-assigning-reviewers-for-Pull-Requests/internal/server"

func main() {
	// config
	// logger
	srv := server.NewHTTPServer(":8080")
	srv.Run()
	// logic
}
