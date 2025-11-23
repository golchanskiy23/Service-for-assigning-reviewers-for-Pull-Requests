package server

import (
	"Service-for-assigning-reviewers-for-Pull-Requests/config"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = "0.0.0.0:3333"
	defaultShutdownTimeout = 10 * time.Second
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
		channelErr:      make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}
	for _, option := range options {
		option(server)
	}
	server.Start()
	return server
}

func (s *Server) FullShutdownTimeout() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	log.Println("Shutting down server...")
	if err := s.internalServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown filed: %v", err)
	}
	log.Println("Server shutdown successfully")
	return nil
}

func (s *Server) GracefulShutdown(logger *slog.Logger) error {
	osInterruptChan := make(chan os.Signal, 1)
	signal.Notify(osInterruptChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-osInterruptChan:
		logger.Info("server interrupted by system or user")
	case err := <-s.channelErr:
		logger.Error("server threw an error", slog.Any("error", err))
	}

	close(osInterruptChan)
	if err := s.FullShutdownTimeout(); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
		return fmt.Errorf("graceful shutdown collapsed: %w", err)
	}

	logger.Info("server shutdown completed successfully")
	return nil
}

func StartServer(cfg *config.Config, controller *chi.Mux, logger *slog.Logger) error {
	customServer := NewServer(controller,
		SetReadTimeout(6*time.Second),
		SetWriteTimeout(6*time.Second),
		SetAddr(),
		SetShutdownTimeout(cfg.Server.ShutdownTimeout),
	)
	logger.Info("successfully created server\n")
	if err := customServer.GracefulShutdown(logger); err != nil {
		return fmt.Errorf("server shutdown error: %v", err)
	}
	return nil
}
