package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterPingRoutes groups lightweight utility endpoints.
// Add new utility endpoints here for consistent structure.
func RegisterPingRoutes(router gin.IRoutes) {
	router.GET("/ping", pingHandler())
}

func pingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"time":    time.Now().UTC().Format(time.RFC3339),
			"userID": userID,
		})
	}
}
