package repo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func RegisterRepoRoutes(router gin.IRoutes, db *pgxpool.Pool) {
	router.GET("/repo", GetRepoHandler(db))
}

func GetRepoHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		name := c.Query("name")

		query := `
			SELECT id, name, author, description, created_at
			FROM repositories
		`
		args := []any{}

		if name != "" {
			query += " WHERE name = ?"
			args = append(args, name)
		}

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		repos := []Repo{}

		for rows.Next() {
			var r Repo
			if err := rows.Scan(
				&r.ID,
				&r.Name,
				&r.Author,
				&r.Description,
				&r.CreatedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			repos = append(repos, r)
		}

		c.JSON(http.StatusOK, repos)
	}
}
