// This is the bridge between the Go code and the database.
// It wraps sql-generated queries in a clean struct your other packages can use
// Theres an important concept here: the connection pool. When your scheduler runs
// and fires 3 goroutines simultaneously (stripe, Gmail, Calendly), all 3 might need
// the database at the same time. A connection pool maintains multiple open connections
// to postgres and lends them out as needed. pgxpool handles this automatically

package storage

import (
	"context"
	"fmt"

	"github.com/Justified02/abm/internal/storage/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store wraps the database connection pool and the sqlc-generated queries struct
// Every database operation in this app goes through this type
type Store struct {
	pool 	*pgxpool.Pool
	queries *db.Queries
}

// New creates a connection pool and verifies the database is reachable.
// This is called once at startup
func New(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	// Ping verifies the credentials and network path are correct
	// Without this, a wrong DATABASE_URL would only fail when the first query runs
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &Store{
		pool: 	pool,
		queries: db.New(pool),
	},  nil
}


// Close releases all connections in the pool. Called via defer in main.go when the program shuts down
func (s *Store) Close() {
	s.pool.Close()
}

// Queries returns the sqlc query interface. Other packages use this to run db operations
func (s *Store) Queries() *db.Queries {
	return s.queries
}