package dbmigrate

import (
	"context"
	"errors"
	"io/fs"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	"questmaster-core/internal/shared/testdb"
	"questmaster-core/migrations"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

// createDatabase creates an empty database next to TEST_DB_URL, dropped when the test ends,
// so these tests never touch the database the repository tests use.
func createDatabase(t *testing.T, name string) string {
	t.Helper()
	admin := testdb.Pool(t)
	ctx := context.Background()

	if _, err := admin.Exec(ctx, "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)"); err != nil {
		t.Fatalf("dropping %s: %v", name, err)
	}
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		t.Fatalf("creating %s: %v", name, err)
	}
	t.Cleanup(func() {
		_, _ = admin.Exec(ctx, "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
	})

	u, err := url.Parse(os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatalf("parsing TEST_DB_URL: %v", err)
	}
	u.Path = "/" + name
	return u.String()
}

// dropVersionTable makes a database look like one migrated by hand, with no record of applied versions
func dropVersionTable(t *testing.T, dbURL string) {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("connecting to %s: %v", dbURL, err)
	}
	defer pool.Close()

	if _, err := pool.Exec(context.Background(), "DROP TABLE schema_migrations"); err != nil {
		t.Fatalf("dropping schema_migrations: %v", err)
	}
}

// latestVersion is the highest migration number in the embedded files
func latestVersion(t *testing.T) uint {
	t.Helper()
	files, err := fs.Glob(migrations.FS, "*.up.sql")
	if err != nil || len(files) == 0 {
		t.Fatalf("listing migrations: %v", err)
	}
	var latest uint
	for _, f := range files {
		v, err := strconv.ParseUint(strings.SplitN(f, "_", 2)[0], 10, 32)
		if err != nil {
			t.Fatalf("migration %s has no numeric version", f)
		}
		latest = max(latest, uint(v))
	}
	return latest
}

func TestMigrations(t *testing.T) {
	dbURL := createDatabase(t, "questmaster_migrate_test")
	latest := latestVersion(t)

	t.Run("empty database gets every migration", func(t *testing.T) {
		version, err := Up(dbURL)
		if err != nil || version != latest {
			t.Fatalf("expected version %d, got %d err=%v", latest, version, err)
		}
	})

	t.Run("up to date database is a no-op", func(t *testing.T) {
		version, err := Up(dbURL)
		if err != nil || version != latest {
			t.Fatalf("expected version %d, got %d err=%v", latest, version, err)
		}
	})

	t.Run("full rollback and reapply", func(t *testing.T) {
		m, err := newMigrate(migrations.FS, dbURL)
		if err != nil {
			t.Fatalf("new migrate: %v", err)
		}
		downErr := m.Down()
		m.Close()
		if downErr != nil {
			t.Fatalf("rolling back: %v", downErr)
		}

		version, err := Up(dbURL)
		if err != nil || version != latest {
			t.Fatalf("expected version %d after reapplying, got %d err=%v", latest, version, err)
		}
	})

	t.Run("baseline a database migrated by hand", func(t *testing.T) {
		// Recreate the production situation: 0001-0004 applied, but no version table
		m, err := newMigrate(migrations.FS, dbURL)
		if err != nil {
			t.Fatalf("new migrate: %v", err)
		}
		migrateErr := m.Migrate(4)
		m.Close()
		if migrateErr != nil {
			t.Fatalf("migrating to 4: %v", migrateErr)
		}
		dropVersionTable(t, dbURL)

		if err := Force(dbURL, 4); err != nil {
			t.Fatalf("forcing version 4: %v", err)
		}
		version, err := Up(dbURL)
		if err != nil || version != latest {
			t.Fatalf("expected version %d after baseline, got %d err=%v", latest, version, err)
		}
	})
}

func TestFailedMigrationStopsLaterRuns(t *testing.T) {
	dbURL := createDatabase(t, "questmaster_migrate_fail_test")
	broken := fstest.MapFS{
		"0001_broken.up.sql":   {Data: []byte("SELECT * FROM table_that_does_not_exist;")},
		"0001_broken.down.sql": {Data: []byte("SELECT 1;")},
	}

	if _, err := up(broken, dbURL); err == nil {
		t.Fatalf("expected the broken migration to fail")
	}

	// The database is now marked dirty, so a new run must refuse to continue
	_, err := up(broken, dbURL)
	var dirty migrate.ErrDirty
	if !errors.As(err, &dirty) {
		t.Fatalf("expected ErrDirty on the next run, got %v", err)
	}
}
