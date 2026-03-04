package repo

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/greedypanda0/kuro/api/remote/database"
	coredb "github.com/greedypanda0/kuro/core/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSnapshotsRoutes(router gin.IRoutes, db *pgxpool.Pool) {
	router.GET("repositories/:id/snapshots", getSnapshots(db))
	router.GET("repositories/:id/snapshots/:snapshot_id", getSnapshot(db))
	router.GET("repositories/:id/snapshots/:snapshot_id/files", getSnapshotFiles(db))
	router.GET("repositories/:id/snapshots/:snapshot_id/files/*file_id", getSnapshotFile(db))
}

func getSnapshots(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		snapshots, err := coredb.ListSnapshots(coredbConnection)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"snapshots": snapshots})
	}
}

func getSnapshot(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		snapshotID := c.Param("snapshot_id")
		snapshot, err := coredb.GetSnapshot(coredbConnection, snapshotID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"snapshot": snapshot})
	}
}

func getSnapshotFiles(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		repoID := c.Param("id")
		snapshotID := c.Param("snapshot_id")

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

		files, err := coredb.ListSnapshotFiles(coredbConnection, snapshotID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"files": files})
	}
}

func getSnapshotFile(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		repoID := c.Param("id")
		snapshotID := c.Param("snapshot_id")
		fileID := c.Param("file_id")
		if len(fileID) > 0 {
			fileID = fileID[1:]
		}

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

		file, err := coredb.GetSnapshotFile(coredbConnection, snapshotID, fileID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"file": file})
	}
}
