package repo

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/greedypanda0/kuro/api/remote/database"
	coredb "github.com/greedypanda0/kuro/core/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRefsRoutes(router gin.IRoutes, db *pgxpool.Pool) {
	router.GET("repositories/:id/refs", getRefs(db))
	router.GET("repositories/:id/refs/:ref", getRef(db))
}

func getRefs(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// userID := c.MustGet("user_id")
		repoID := c.Param("id")
		repo, err := database.GetRepo(db, c, repoID)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		if repo == nil {
			c.JSON(404, gin.H{"error": "repository not found"})
			return
		}

		path := filepath.Join("data", repo.UserID, repo.Name+".db")

		coredbConnection, err := coredb.OpenDB(path)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer coredbConnection.Close()

		refs, err := coredb.ListRefs(coredbConnection)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"refs": refs})
	}
}

func getRef(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		repoID := c.Param("id")
		refName := c.Param("ref")
		repo, err := database.GetRepo(db, c, repoID)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		if repo == nil {
			c.JSON(404, gin.H{"error": "repository not found"})
			return
		}

		path := filepath.Join("data", repo.UserID, repo.Name+".db")

		coredbConnection, err := coredb.OpenDB(path)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer coredbConnection.Close()

		ref, err := coredb.GetRef(coredbConnection, refName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"ref": ref})
	}
}
