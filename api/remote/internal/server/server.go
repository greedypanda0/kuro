package server

import (
	"net/http"
	"time"

	"api/remote/database"
	"api/remote/internal/config"
	"api/remote/internal/handlers"
	"api/remote/internal/logger"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config, log *logger.Logger) *http.Server {
	if cfg.Log.Development {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err.Error())
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
