package logger

import (
	"context"
	"log/slog"
	"os"
)

type MultiLeveLHandler struct {
	infoHandler  slog.Handler
	debugHandler slog.Handler
	errorHandler slog.Handler
}

func (MultiLeveLHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

//nolint:gocritic // hugeParam: slog.Record is part of standard interface
func (handler MultiLeveLHandler) Handle(
	ctx context.Context,
	r slog.Record,
) error {
	switch r.Level {
	case slog.LevelInfo:
		return handler.infoHandler.Handle(ctx, r)
	case slog.LevelDebug:
		return handler.debugHandler.Handle(ctx, r)
	case slog.LevelError:
		return handler.errorHandler.Handle(ctx, r)
	case slog.LevelWarn:
		return handler.infoHandler.Handle(ctx, r)
	default:
		return handler.infoHandler.Handle(ctx, r)
	}
}

func (handler MultiLeveLHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MultiLeveLHandler{
		infoHandler:  handler.infoHandler.WithAttrs(attrs),
		debugHandler: handler.debugHandler.WithAttrs(attrs),
		errorHandler: handler.errorHandler.WithAttrs(attrs),
	}
}

func (handler MultiLeveLHandler) WithGroup(name string) slog.Handler {
	return &MultiLeveLHandler{
		infoHandler:  handler.infoHandler.WithGroup(name),
		debugHandler: handler.debugHandler.WithGroup(name),
		errorHandler: handler.errorHandler.WithGroup(name),
	}
}

func SetupLogger() *slog.Logger {
	handler := &MultiLeveLHandler{
		infoHandler: slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		),
		debugHandler: slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
		errorHandler: slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelError},
		),
	}

	return slog.New(handler)
}
