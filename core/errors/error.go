package errors

import "errors"

var (
	ErrRepoAlreadyInitialized = errors.New("repository already initialized")
	ErrRepoNotInitialized     = errors.New("repository not initialized")
	ErrDatabaseOpenFailed     = errors.New("failed to open database")
	ErrDatabasePingFailed     = errors.New("failed to ping database")
	ErrSchemaApplyFailed      = errors.New("failed to apply schema")
	ErrDataNotFound           = errors.New("data not found")
)
