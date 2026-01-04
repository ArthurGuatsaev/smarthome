package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArthurGuatsaev/smarthome/internal/app"
	"github.com/ArthurGuatsaev/smarthome/internal/config"
	"github.com/ArthurGuatsaev/smarthome/internal/httpapi"
	"github.com/ArthurGuatsaev/smarthome/internal/storage"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	db, err := storage.Open(context.Background(), cfg.DBPath)
	if err != nil {
		slog.Error("db_open_error", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	migCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := storage.Migrate(migCtx, db.DB); err != nil {
		slog.Error("db_migrate_error", "err", err)
		os.Exit(1)
	}

	application := app.New(db)
	srv := httpapi.NewServer(application, cfg.APIKey)
	httpServer := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      srv.Handler(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server_start", "addr", cfg.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server_error", "err", err)
			stop()
		}
	}()

	<-ctx.Done()

	slog.Info("server_shutdown")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown_error", "err", err)
	} else {
		slog.Info("shutdown_ok")
	}
}

func setupLogger(level string) {
	lvl := slog.LevelInfo
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	slog.SetDefault(slog.New(handler))
}
