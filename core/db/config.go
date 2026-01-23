package db

import (
	"core/errors"
	"database/sql"
)

func SetConfig(db *sql.DB, key, value string) error {
	_, err := db.Exec(
		"INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)",
		key,
		value,
	)
	return err
}

func GetConfig(db *sql.DB, key string) (string, error) {
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

func DeleteConfig(db *sql.DB, key string) error {
	_, err := db.Exec(
		"DELETE FROM config WHERE key = ?",
		key,
	)
	return err
}
