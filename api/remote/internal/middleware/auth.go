package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")
		token, isBearer := strings.CutPrefix(authHeader, "Bearer ")
		if isBearer {
			var userID string

			err := pool.QueryRow(ctx, `
				SELECT user_id
				FROM auth_tokens
				WHERE id = $1
				AND (expires_at IS NULL OR expires_at > NOW())
			`, token).Scan(&userID)

			if err != nil {
				if err == pgx.ErrNoRows {
					c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
					return
				}
				c.AbortWithStatusJSON(500, gin.H{"error": "db error"})
				return
			}

			c.Set("user_id", userID)
			c.Next()
			return
		}

		sessionToken, err := c.Cookie("authjs.session-token")
		if err != nil {
			sessionToken, err = c.Cookie("__Secure-authjs.session-token")
			if err != nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
				return
			}
		}

		var userID string
		err = pool.QueryRow(ctx, `
			SELECT user_id
			FROM sessions
			WHERE session_token = $1
			AND expires > NOW()
		`, sessionToken).Scan(&userID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired session"})
				return
			}
			c.AbortWithStatusJSON(500, gin.H{"error": "db error"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
