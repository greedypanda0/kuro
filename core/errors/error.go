package errors

import "errors"

var (
	ErrRepoAlreadyInitialized = errors.New("repository already initialized")
	ErrRepoNotInitialized     = errors.New("repository not initialized")
	ErrDatabaseOpenFailed     = errors.New("failed to open database")
	ErrDatabasePingFailed     = errors.New("failed to ping database")
	ErrSchemaApplyFailed      = errors.New("failed to apply schema")
	ErrDataNotFound           = errors.New("data not found")
	ErrRefNotFound            = errors.New("ref not found")
	ErrSnapshotNotFound       = errors.New("snapshot not found")
	ErrObjectNotFound         = errors.New("object not found")
	ErrIgnoreFileNotFound     = errors.New("ignore file not found")
)
