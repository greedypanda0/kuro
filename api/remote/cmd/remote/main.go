package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/greedypanda0/kuro/api/remote/database"
	"github.com/greedypanda0/kuro/api/remote/internal/config"
	"github.com/greedypanda0/kuro/api/remote/internal/logger"
	"github.com/greedypanda0/kuro/api/remote/internal/server"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	log := logger.Init(cfg.Log)
	defer func() {
		_ = log.Sync()
	}()

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	httpServer := server.New(cfg, log, db)

	go func() {
		log.Info("remote api starting", logger.String("addr", cfg.HTTP.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server crashed", logger.Error(err))
		}
	}()

	shutdownOnSignal(log, httpServer, cfg.HTTP.ShutdownTimeout, db)
}

func shutdownOnSignal(log *logger.Logger, httpServer *http.Server, timeout time.Duration, db *pgxpool.Pool) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Info("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	defer db.Close()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", logger.Error(err))
		if err := httpServer.Close(); err != nil {
			log.Error("forced shutdown failed", logger.Error(err))
		}
		return
	}

	log.Info("shutdown complete")
}
