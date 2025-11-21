package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Service-for-assigning-reviewers-for-Pull-Requests/config"

	"github.com/go-chi/chi/v5"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = "0.0.0.0:3333"
	defaultShutdownTimeout = 10 * time.Second
	chanBufferSize         = 1
	Billion                = 1000000000
)

type Server struct {
	internalServer  *http.Server
	channelErr      chan error
	shutdownTimeout time.Duration
}

func (s *Server) Start() {
	go func() {
		s.channelErr <- s.internalServer.ListenAndServe()
		close(s.channelErr)
	}()
}

func NewServer(handler http.Handler, options ...Option) *Server {
	server := &Server{
		internalServer: &http.Server{
			Handler:      handler,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			Addr:         defaultAddr,
		},
		channelErr:      make(chan error, chanBufferSize),
		shutdownTimeout: defaultShutdownTimeout,
	}
	for _, option := range options {
		option(server)
	}

	server.Start()

	return server
}

func (s *Server) FullShutdownTimeout(logger *slog.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	logger.Info("Shutting down server...\n")

	if err := s.internalServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown filed: %w", err)
	}

	return nil
}

func (s *Server) GracefulShutdown(logger *slog.Logger) {
	osInterruptChan := make(chan os.Signal, chanBufferSize)
	signal.Notify(osInterruptChan, syscall.SIGTERM, syscall.SIGINT)

	timeoutChan := time.After(s.shutdownTimeout)

	select {
	case <-osInterruptChan:
		logger.Info("server interrupted by system or user")
	case <-s.channelErr:
		logger.Error("server error occurred", slog.Any("error", <-s.channelErr))
	case <-timeoutChan:
		logger.Info("shutdown timeout reached")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	logger.Info("shutting down server...")

	if err := s.internalServer.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
	} else {
		logger.Info("server stopped successfully")
	}

	close(osInterruptChan)
}

func StartServer(
	cfg *config.Config,
	controller *chi.Mux,
	logger *slog.Logger,
) *Server {
	customServer := NewServer(controller,
		SetReadTimeout(*cfg.Server.ReadTimeout),
		SetWriteTimeout(*cfg.Server.WriteTimeout),
		SetAddr(cfg.Server.Addr),
		SetShutdownTimeout(cfg.Server.ShutdownTimeout),
	)
	logger.Info("server shutdown info",
		"shutdownTimeout", fmt.Sprintf("%d %s",
			customServer.shutdownTimeout/Billion, "s"),
	)
	logger.Info("")
	logger.Info("successfully created server\n")

	return customServer
}
