package repositories

import "database/sql"

type DatabaseService interface {
	Connect() error
	Close() error
	Health() map[string]string

	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row

	BeginTransaction() (*sql.Tx, error)
}
