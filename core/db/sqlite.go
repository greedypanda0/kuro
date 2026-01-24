package db

import (
	"core/errors"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitSQL(path string) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create repo dir: %w", err)
	}

	databasePath := path

	if _, err := os.Stat(databasePath); err == nil {
		return nil, errors.ErrRepoAlreadyInitialized
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("stat db: %w", err)
	}

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOpenFailed, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabasePingFailed, err)
	}

	return db, nil
}

func OpenDB(path string) (*sql.DB, error) {
	databasePath := path

	if _, err := os.Stat(databasePath); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.ErrRepoNotInitialized
		}
		return nil, fmt.Errorf("stat db: %w", err)
	}

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabaseOpenFailed, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: %v", errors.ErrDatabasePingFailed, err)
	}

	return db, nil
}

func ApplySchema(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("%w: %v", errors.ErrSchemaApplyFailed, err)
	}
	return nil
}
