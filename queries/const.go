package queries

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound            = sql.ErrNoRows
	ErrForeignKeyViolation = errors.New("foreign key violation")
)
