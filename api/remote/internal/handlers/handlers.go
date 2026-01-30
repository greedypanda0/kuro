package handlers

import (
	"net/http"
	"time"

	"api/remote/internal/build"
	"api/remote/internal/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, log *logger.Logger) {
	router.Use(requestLogger(log))
	router.Use(gin.Recovery())

	router.GET("/health", healthHandler())
	router.GET("/version", versionHandler())
	RegisterPingRoutes(router)
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
