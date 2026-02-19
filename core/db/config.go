package db

import (
	"github.com/greedypanda0/kuro/core/errors"
	"database/sql"
)

type Config struct {
	Key   string
	Value string
}

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

func ListConfigs(db DBTX) ([]Config, error) {
	rows, err := db.Query("SELECT key, value FROM config")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []Config
	for rows.Next() {
		var config Config
		if err := rows.Scan(&config.Key, &config.Value); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}
