package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"regexp"
	"time"
)

const URLTemplate = "postgresql://%s:%s@%s:%s/%s?sslmode=disable" // where %s and %d defined as: DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME

var (
	ErrNotFound  = fmt.Errorf("not found")
	ErrDuplicate = fmt.Errorf("duplicate")
)

type IDB interface {
	bun.IDB
	ID() string // id for the database, used for transactions
}

type DB struct {
	*bun.DB
	id string
}

func (db *DB) BunDB() *bun.DB {
	return db.DB
}

func (db *DB) SqlDB() *sql.DB {
	return db.DB.DB
}

func (db *DB) ID() string {
	return db.id
}

type Tx struct {
	bun.Tx
	id string
}

func (t Tx) ID() string {
	return t.id
}

// NewDatabase creates new SQL database instance.
func NewDatabase(dsn string, debug bool) (*DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	bundb := bun.NewDB(sqldb, pgdialect.New())
	// log queries when debug mode is set
	bundb.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithEnabled(debug),
			bundebug.WithVerbose(true),
		),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := bundb.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	db := &DB{DB: bundb, id: uuid.Must(uuid.NewUUID()).String()}
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
	} else if matched, _ := regexp.Match(`^ERROR: duplicate key value violates unique constraint .* \(SQLSTATE=23505\)$`, []byte(err.Error())); matched {
		return true
	}
	return false
}

func MigrateDBSchema(db *sql.DB, migrationsPath, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		dbName,
		driver)
	if err != nil {
		return err
	}

	return m.Up()
}

type TransactionFunc func(ctx context.Context, f func(tctx context.Context) error) error

func NewTransactionFunc(db IDB) TransactionFunc {
	return func(ctx context.Context, f func(tctx context.Context) error) error {
		ff := func(ctx context.Context, tx bun.Tx) error {
			tctx := context.WithValue(ctx, txKey{}, &Tx{Tx: tx, id: db.ID()})
			return f(tctx)
		}
		return db.RunInTx(ctx, nil, ff)
	}
}

// DbTx returns transaction or uses current db
func DbTx(db IDB, ctx context.Context) IDB {
	if tx, ok := ctx.Value(txKey{}).(IDB); ok && tx.ID() == db.ID() {
		return tx
	}
	return db
}
