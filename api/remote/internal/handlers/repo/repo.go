package repo

import (
	"encoding/json"
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
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type uploadMetadata struct {
	Name string `json:"name"`
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
		author := c.Query("author")

		query := `
			SELECT id, name, author, description, created_at
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

		if author != "" {
			conditions = append(conditions, fmt.Sprintf("author = $%d", i))
			args = append(args, author)
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
				&r.Author,
				&r.Description,
				&r.CreatedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to scan repository",
				})
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
			SELECT id, name, author, description, created_at
			FROM repositories
			WHERE id = $1
		`

		row := db.QueryRow(ctx, query, id)

		var r Repository
		err := row.Scan(
			&r.ID,
			&r.Name,
			&r.Author,
			&r.Description,
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

		var metadata uploadMetadata
		var fileReader io.Reader
		userID := c.MustGet("user_id").(string)

		contentType := c.ContentType()

		switch {
		case strings.HasPrefix(contentType, "multipart/form-data"):
			if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
				return
			}

			metadataStr := c.PostForm("metadata")
			if metadataStr == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing metadata"})
				return
			}

			if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metadata"})
				return
			}

			fileHeader, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
				return
			}

			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
				return
			}
			defer file.Close()

			fileReader = file

		case strings.HasPrefix(contentType, "application/octet-stream"):
			metadata.Name = c.GetHeader("X-Name")

			if metadata.Name == "" {
				metadata.Name = c.Query("name")
			}

			fileReader = c.Request.Body

		default:
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "unsupported content type"})
			return
		}

		if userID == "" || metadata.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id or db_name"})
			return
		}

		dirPath := filepath.Join("data", userID)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user data directory"})
			return
		}

		dstPath := filepath.Join(dirPath, "temp_"+metadata.Name+".db")
		dstFile, err := os.Create(dstPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create database file"})
			return
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, fileReader); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save database file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"path": dstPath,
		})
	}
}
