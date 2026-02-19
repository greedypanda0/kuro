package db

import (
	"github.com/greedypanda0/kuro/core/errors"
	"database/sql"
)

type Object struct {
	Hash      string
	Content   []byte
	CreatedAt int64
}

func CreateObject(db DBTX, hash string, content []byte) error {
	_, err := db.Exec(
		"INSERT OR IGNORE INTO objects (hash, content) VALUES (?, ?)",
		hash,
		content,
	)
	return err
}

func GetObject(db DBTX, hash string) (*Object, error) {
	var obj Object

	err := db.QueryRow(
		"SELECT hash, content, created_at FROM objects WHERE hash = ?",
		hash,
	).Scan(&obj.Hash, &obj.Content, &obj.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.ErrObjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func DeleteObject(db DBTX, hash string) error {
	res, err := db.Exec(
		"DELETE FROM objects WHERE hash = ?",
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
		return errors.ErrObjectNotFound
	}

	return nil
}

func ListObjects(db DBTX) ([]Object, error) {
	rows, err := db.Query("SELECT hash, content, created_at FROM objects ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objects []Object
	for rows.Next() {
		var obj Object
		if err := rows.Scan(&obj.Hash, &obj.Content, &obj.CreatedAt); err != nil {
			return nil, err
		}
		objects = append(objects, obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}
