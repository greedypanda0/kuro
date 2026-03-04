package repo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/greedypanda0/kuro/api/remote/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

		var repos []database.Repository

		for rows.Next() {
			var r database.Repository
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
		id := c.Param("id")
		repo, err := database.GetRepo(db, c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if repo == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "repo not found",
			})
			return
		}

		c.JSON(http.StatusOK, repo)
	}
}

func postRepositoryHandler(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

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
		repoName := remoteParts[1]
		repo, err := database.GetRepo(db, c, repoName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if repo == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "repo not found",
			})
			return
		}

		if repo.UserID != userID {
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

		dstPath := filepath.Join(dirPath, "temp_"+repo.Name+".db")

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

		finalPath := filepath.Join(dirPath, repo.Name+".db")
		tempPath := filepath.Join(dirPath, "temp_"+repo.Name+".db")

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
