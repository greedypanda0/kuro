package repo

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func RegisterRepositoryRoutes(router gin.IRoutes, db *pgxpool.Pool) {
	router.GET("/repositories", getRepositoriesHandler(db))
	router.GET("/repositories/:id", getRepositoryHandler(db))
	router.POST("/repositories", postRepositoryHandler(db))
}

func getRepositoriesHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		name := c.Query("name")
		userID := c.Query("user_id")

		query := `
			SELECT id, name, user_id, created_at
			FROM repositories
		`

		var args []any
		var conditions []string
		i := 1

		if name != "" {
			conditions = append(conditions, fmt.Sprintf("name = $%d", i))
			args = append(args, name)
			i++
		}

		if userID != "" {
			conditions = append(conditions, fmt.Sprintf("user_id = $%d", i))
			args = append(args, userID)
			i++
		}

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch repositories",
			})
			return
		}
		defer rows.Close()

		var repos []Repository

		for rows.Next() {
			var r Repository
			if err := rows.Scan(
				&r.ID,
				&r.Name,
				&r.UserID,
				&r.CreatedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to scan repository",
				})
				fmt.Println("Error", err.Error())
				return
			}
			repos = append(repos, r)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "row iteration error",
			})
			return
		}

		c.JSON(http.StatusOK, repos)
	}
}

func getRepositoryHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		id := c.Param("id")

		query := `
			SELECT id, name, user_id, created_at
			FROM repositories
			WHERE id = $1
		`

		row := db.QueryRow(ctx, query, id)

		var r Repository
		err := row.Scan(
			&r.ID,
			&r.Name,
			&r.UserID,
			&r.CreatedAt,
		)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "repository not found",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch repository",
			})
			return
		}

		c.JSON(http.StatusOK, r)
	}
}

func postRepositoryHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)
		ctx := c.Request.Context()

		if !strings.HasPrefix(c.ContentType(), "application/octet-stream") {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error": "content type must be application/octet-stream",
			})
			return
		}

		remote := c.GetHeader("X-Remote")
		if remote == "" {
			remote = c.Query("remote")
		}

		if remote == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing remote",
			})
			return
		}

		remoteParts := strings.SplitN(remote, "/", 2)
		// user := remoteParts[0]
		repo := remoteParts[1]
		row := db.QueryRow(ctx, "SELECT user_id FROM repositories WHERE name = $1", repo)
		var repoUserID string
		if err := row.Scan(&repoUserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch user",
			})
			return
		}
		
		if repoUserID != userID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "repository does not belong to user",
			})
			return
		}

		dirPath := filepath.Join("data", userID)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create user data directory",
			})
			return
		}

		dstPath := filepath.Join(dirPath, "temp_"+repo+".db")

		dstFile, err := os.Create(dstPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create database file",
			})
			return
		}

		if _, err := io.Copy(dstFile, c.Request.Body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save database file",
			})
			return
		}

		if err := dstFile.Sync(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to sync database file",
			})
			return
		}

		if err := dstFile.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to close database file",
			})
			return
		}

		finalPath := filepath.Join(dirPath, repo+".db")
		tempPath := filepath.Join(dirPath, "temp_"+repo+".db")

		if err := os.Remove(finalPath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to remove old repository",
			})
			return
		}

		if err := os.Rename(tempPath, finalPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to finalize repository",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"path": dstPath,
		})
	}
}
