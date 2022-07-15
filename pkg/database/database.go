package database

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"regexp"
	"time"
)

const DbURLTemplate = "postgresql://%s:%s@%s:%d/%s?sslmode=disable" // where %s and %d defined as: DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME

// NewDatabase creates new SQL database instance.
func NewDatabase(dsn string, debug bool) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	// log queries when debug mode is set
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithEnabled(debug),
			bundebug.WithVerbose(true),
		),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// IsNotFound check if database didn't return any rows.
func IsNotFound(err error) bool {
	if err != nil && err.Error() == "sql: no rows in result set" {
		return true
	}
	return false
}

// IsDuplicate checks if an entity already exists in the database.
// Note: This regexp pattern is for postgres. Add more patterns for other databases.
func IsDuplicate(err error) bool {
	if err == nil {
		return false
	} else if matched, _ := regexp.Match("^ERROR: duplicate key value violates unique constraint .* \\(SQLSTATE=23505\\)$", []byte(err.Error())); matched {
		return true
	}
	return false
}
