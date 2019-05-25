package postgres

import (
	"context"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // package is specific to postgres, so embeds the postgres driver
	"github.com/rs/zerolog"
)

// DB is a database connection handle.
type DB struct {
	*sqlx.DB
	log zerolog.Logger
}

// NewDB returns a new database connection handle.
func NewDB(dsn string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Create a new common logger
	l := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.ErrorLevel)

	return &DB{DB: db, log: l}, nil
}
