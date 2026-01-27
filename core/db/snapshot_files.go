package db

import (
	"core/errors"
	"database/sql"
)

type SnapshotFile struct {
	SnapshotHash string
	Path         string
	ObjectHash   string
}

func CreateSnapshotFile(db DBTX, snapshotHash, filePath, objectHash string) error {
	_, err := db.Exec(
		"INSERT INTO snapshot_files (snapshot_hash, path, object_hash) VALUES (?, ?, ?)",
		snapshotHash,
		filePath,
		objectHash,
	)
	return err
}

func GetSnapshotFile(db DBTX, snapshotHash, filePath string) (*SnapshotFile, error) {
	var sf SnapshotFile
	err := db.QueryRow(
		"SELECT snapshot_hash, path, object_hash FROM snapshot_files WHERE snapshot_hash = ? AND path = ?",
		snapshotHash,
		filePath,
	).Scan(&sf.SnapshotHash, &sf.Path, &sf.ObjectHash)

	if err == sql.ErrNoRows {
		return nil, errors.ErrDataNotFound
	}
	if err != nil {
		return nil, err
	}

	return &sf, nil
}

func DeleteSnapshotFile(db DBTX, snapshotHash, filePath string) error {
	res, err := db.Exec(
		"DELETE FROM snapshot_files WHERE snapshot_hash = ? AND path = ?",
		snapshotHash,
		filePath,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.ErrDataNotFound
	}

	return nil
}

func ListSnapshotFiles(db DBTX, snapshotHash string) ([]SnapshotFile, error) {
	rows, err := db.Query(
		"SELECT snapshot_hash, path, object_hash FROM snapshot_files WHERE snapshot_hash = ? ORDER BY path",
		snapshotHash,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []SnapshotFile
	for rows.Next() {
		var sf SnapshotFile
		if err := rows.Scan(&sf.SnapshotHash, &sf.Path, &sf.ObjectHash); err != nil {
			return nil, err
		}
		files = append(files, sf)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
