package db

import (
	"github.com/greedypanda0/kuro/core/errors"
	"database/sql"
)

type Snapshot struct {
	Hash       string
	ParentHash *string
	Message    string
	Author     *string
	Timestamp  int64
}

func CreateSnapshot(db DBTX, hash string, parentHash *string, message string, author *string) error {
	_, err := db.Exec(
		"INSERT INTO snapshot (hash, parent_hash, message, author) VALUES (?, ?, ?, ?)",
		hash,
		parentHash,
		message,
		author,
	)
	return err
}

func GetSnapshot(db DBTX, hash string) (*Snapshot, error) {
	var (
		s      Snapshot
		parent sql.NullString
		author sql.NullString
	)

	err := db.QueryRow(
		"SELECT hash, parent_hash, message, author, timestamp FROM snapshot WHERE hash = ?",
		hash,
	).Scan(&s.Hash, &parent, &s.Message, &author, &s.Timestamp)

	if err == sql.ErrNoRows {
		return nil, errors.ErrSnapshotNotFound
	}
	if err != nil {
		return nil, err
	}

	if parent.Valid {
		s.ParentHash = &parent.String
	}
	if author.Valid {
		s.Author = &author.String
	}

	return &s, nil
}

func DeleteSnapshot(db DBTX, hash string) error {
	res, err := db.Exec(
		"DELETE FROM snapshot WHERE hash = ?",
		hash,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.ErrSnapshotNotFound
	}

	return nil
}

func ListSnapshots(db DBTX) ([]Snapshot, error) {
	rows, err := db.Query("SELECT hash, parent_hash, message, author, timestamp FROM snapshot ORDER BY timestamp")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []Snapshot
	for rows.Next() {
		var (
			s      Snapshot
			parent sql.NullString
			author sql.NullString
		)

		if err := rows.Scan(&s.Hash, &parent, &s.Message, &author, &s.Timestamp); err != nil {
			return nil, err
		}

		if parent.Valid {
			s.ParentHash = &parent.String
		}
		if author.Valid {
			s.Author = &author.String
		}

		snapshots = append(snapshots, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snapshots, nil
}
