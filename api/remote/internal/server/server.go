package server

import (
	"net/http"
	"time"

	"github.com/greedypanda0/kuro/api/remote/internal/config"
	"github.com/greedypanda0/kuro/api/remote/internal/handlers"
	"github.com/greedypanda0/kuro/api/remote/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg config.Config, log *logger.Logger, db *pgxpool.Pool) *http.Server {
	if cfg.Log.Development {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	handlers.RegisterRoutes(router, log, db)

	return &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
