package db

import (
	"github.com/greedypanda0/kuro/core/errors"
	"database/sql"
)

func SetConfig(db DBTX, key, value string) error {
	_, err := db.Exec(
		"INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)",
		key,
		value,
	)
	return err
}

func GetConfig(db DBTX, key string) (string, error) {
	var value string
	err := db.QueryRow(
		"SELECT value FROM config WHERE key = ?",
		key,
	).Scan(&value)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.ErrDataNotFound
		}
		return "", err
	}

	return value, nil
}

func DeleteConfig(db DBTX, key string) error {
	_, err := db.Exec(
		"DELETE FROM config WHERE key = ?",
		key,
	)
	return err
}
