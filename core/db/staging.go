package db

import (
	"time"
)

type Stage struct {
	Path     string
	StagedAt time.Time
}

func AddStageFile(db DBTX, path string) error {
	_, err := db.Exec(
		"UPSERT INTO staged_files (path) VALUES (?)",
		path,
	)
	return err
}

func RemoveStageFile(db DBTX, path string) error {
	_, err := db.Exec(
		"DELETE FROM staged_files WHERE path = ?",
		path,
	)
	return err
}

func ClearStage(db DBTX) error {
	_, err := db.Exec("DELETE FROM staged_files")
	return err
}

func GetStageFiles(db DBTX) ([]Stage, error) {
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
