package interfaces

import (
	"database/sql"
)

type Driver interface {
	Execute(sql string) (sql.Result, error)
	Query(sql string) (*sql.Rows, error)
	Select(dest interface{}, sql string) error
	Close() error
}
