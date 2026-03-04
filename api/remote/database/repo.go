package database

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func GetRepo(db *pgxpool.Pool, c *gin.Context, repoID string) (*Repository, error) {
	repo := Repository{}
	ctx := c.Request.Context()

	if err := db.QueryRow(ctx, "select id, name, user_id from repositories where id = $1", repoID).Scan(&repo.ID, &repo.Name, &repo.UserID); err != nil {
		return nil, err
	}

	return &repo, nil
}
