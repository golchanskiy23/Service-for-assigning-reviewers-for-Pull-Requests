package server

import (
	"os"
	"time"
)

type Option func(srv *Server)

func SetReadTimeout(duration time.Duration) Option {
	return func(srv *Server) {
		srv.internalServer.ReadTimeout = duration
	}
}

func SetWriteTimeout(duration time.Duration) Option {
	return func(srv *Server) {
		srv.internalServer.WriteTimeout = duration
	}
}

func SetAddr() Option {
	return func(srv *Server) {
		if str := os.Getenv("ADDR"); str != "" {
			srv.internalServer.Addr = str
		}
	}
}

func SetShutdownTimeout(duration time.Duration) Option {
	return func(srv *Server) {
		srv.shutdownTimeout = duration
	}
}
