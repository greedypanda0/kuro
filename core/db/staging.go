package db

import (
	"database/sql"
	"time"
)

type Stage struct {
	Path     string
	StagedAt time.Time
}

func AddStageFile(db *sql.DB, path string) error {
	_, err := db.Exec(
		"INSERT OR IGNORE INTO staged_files (path) VALUES (?)",
		path,
	)
	return err
}

func RemoveStageFile(db *sql.DB, path string) error {
	_, err := db.Exec(
		"DELETE FROM staged_files WHERE path = ?",
		path,
	)
	return err
}

func ClearStage(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM staged_files")
	return err
}

func GetStageFiles(db *sql.DB) ([]Stage, error) {
	rows, err := db.Query(
		"SELECT path, staged_at FROM staged_files ORDER BY staged_at",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []Stage

	for rows.Next() {
		var s Stage
		var ts int64

		if err := rows.Scan(&s.Path, &ts); err != nil {
			return nil, err
		}

		s.StagedAt = time.Unix(ts, 0)
		stages = append(stages, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stages, nil
}
