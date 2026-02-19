package users

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func RegisterUserRoutes(router gin.IRoutes, db *pgxpool.Pool) {
	router.GET("/users/me", getMeRoute(db))
	router.GET("/users", getUsersRoute(db))
}

func getUsersRoute(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		name := c.Query("name")

		pageStr := c.DefaultQuery("page", "1")
		limitStr := c.DefaultQuery("limit", "10")

		page, _ := strconv.Atoi(pageStr)
		limit, _ := strconv.Atoi(limitStr)

		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}

		offset := (page - 1) * limit

		query := `
			SELECT id, name, email
			FROM users
		`

		args := []any{}
		argPos := 1

		if name != "" {
			query += fmt.Sprintf(" WHERE name ILIKE $%d", argPos)
			args = append(args, "%"+name+"%")
			argPos++
		}

		query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argPos, argPos+1)
		args = append(args, limit, offset)

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			users = append(users, u)
		}

		c.JSON(200, users)
	}
}


func getMeRoute(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		ctx := c.Request.Context()
		var user User

		row := db.QueryRow(ctx, "select id, name, email from users where id = $1", userID)
		err := row.Scan(&user.ID, &user.Name, &user.Email)

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch user"})
			return
		}

		c.JSON(200, user)
	}
}
