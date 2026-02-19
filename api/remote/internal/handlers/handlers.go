package handlers

import (
	"net/http"
	"time"

	"api/remote/internal/build"
	"api/remote/internal/handlers/repo"
	"api/remote/internal/handlers/users"
	"api/remote/internal/logger"
	"api/remote/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.Engine, log *logger.Logger, db *pgxpool.Pool) {
	router.Use(requestLogger(log))
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.AuthMiddleware(db))

	apiRouter.GET("/health", healthHandler())
	apiRouter.GET("/version", versionHandler())
	RegisterPingRoutes(apiRouter)
	repo.RegisterRepositoryRoutes(apiRouter, db)
	users.RegisterUserRoutes(apiRouter, db)

}

func healthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func versionHandler() gin.HandlerFunc {
	info := build.Info()
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, info)
	}
}

func requestLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		log.Info("http request",
			logger.String("method", c.Request.Method),
			logger.String("path", c.FullPath()),
			logger.Int("status", c.Writer.Status()),
			logger.Duration("duration", time.Since(start)),
		)
	}
}
