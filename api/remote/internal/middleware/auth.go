package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Session struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func AuthMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("authjs.session-token")
		if err != nil {
			sessionToken, err = c.Cookie("__Secure-authjs.session-token")
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "no session cookie"})
				return
			}
		}

		fmt.Println("SESSION TOKEN:", sessionToken)

		ctx := c.Request.Context()

		var userID string

		err = pool.QueryRow(
			ctx,
			`SELECT user_id
			 FROM sessions
			 WHERE session_token = $1
			   AND expires > now()`,
			sessionToken,
		).Scan(&userID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.AbortWithStatusJSON(401, gin.H{"error": "invalid session"})
				return
			}
			c.AbortWithStatusJSON(500, gin.H{"error": "db error"})
			return
		}

		c.Set("user_id", userID)

		c.Next()
	}
}
