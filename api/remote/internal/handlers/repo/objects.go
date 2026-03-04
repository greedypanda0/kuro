package repo

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/greedypanda0/kuro/api/remote/database"
	coredb "github.com/greedypanda0/kuro/core/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Object struct {
	Hash      string
	CreatedAt int64
}

func RegisterObjectsRoutes(router gin.IRoutes, db *pgxpool.Pool) {
	router.GET("repositories/:id/objects", getObjects(db))
	router.GET("repositories/:id/objects/:hash", getObject(db))
}

func getObjects(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		repoID := c.Param("id")
		repo, err := database.GetRepo(db, c, repoID)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		path := filepath.Join("data", repo.UserID, repo.Name+".db")

		coredbConnection, err := coredb.OpenDB(path)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer coredbConnection.Close()

		var objects []Object
		rawObjects, err := coredb.ListObjects(coredbConnection)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		for _, rawObject := range rawObjects {
			objects = append(objects, Object{
				Hash:      rawObject.Hash,
				CreatedAt: rawObject.CreatedAt,
			})
		}

		c.JSON(200, gin.H{"objects": objects})
	}
}

func getObject(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		repoID := c.Param("id")
		repo, err := database.GetRepo(db, c, repoID)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		path := filepath.Join("data", repo.UserID, repo.Name+".db")

		coredbConnection, err := coredb.OpenDB(path)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer coredbConnection.Close()

		hash := c.Param("hash")
		object, err := coredb.GetObject(coredbConnection, hash)
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"object": object})
	}
}
