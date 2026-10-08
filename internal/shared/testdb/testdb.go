// Package testdb gives repository tests a connection to a real Postgres database.
package testdb

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool connects to the database in TEST_DB_URL, which must have every migration applied.
// It uses its own variable instead of DB_URL so tests never run against a development database by accident.
// The test is skipped when TEST_DB_URL is unset.
func Pool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		t.Skip("TEST_DB_URL not set, skipping database test")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("connecting to test database: %v", err)
	}
	t.Cleanup(pool.Close)

	return pool
}
