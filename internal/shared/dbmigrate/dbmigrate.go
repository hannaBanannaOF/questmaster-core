// Package dbmigrate applies the embedded SQL migrations with golang-migrate.
package dbmigrate

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"questmaster-core/migrations"
)

// Up applies every pending migration and returns the resulting version.
// It is a no-op when the database is already up to date.
func Up(dbURL string) (uint, error) {
	return up(migrations.FS, dbURL)
}

// Force records version as applied and clears a failed (dirty) state, without running any migration.
// Use it to baseline a database migrated by hand, or to recover after fixing a failed migration.
func Force(dbURL string, version int) error {
	m, err := newMigrate(migrations.FS, dbURL)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Force(version)
}

func up(fsys fs.FS, dbURL string) (uint, error) {
	m, err := newMigrate(fsys, dbURL)
	if err != nil {
		return 0, err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, err
	}

	version, _, err := m.Version()
	return version, err
}

func newMigrate(fsys fs.FS, dbURL string) (*migrate.Migrate, error) {
	source, err := iofs.New(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("reading migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, pgx5URL(dbURL))
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return m, nil
}

// pgx5URL rewrites a postgres:// URL to pgx5://, since golang-migrate picks its database driver from the scheme.
func pgx5URL(dbURL string) string {
	for _, scheme := range []string{"postgres://", "postgresql://"} {
		if strings.HasPrefix(dbURL, scheme) {
			return "pgx5://" + strings.TrimPrefix(dbURL, scheme)
		}
	}
	return dbURL
}
