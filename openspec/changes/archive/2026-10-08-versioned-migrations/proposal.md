# Proposal

## Why

Migrations are applied by hand in production, and the database does not record which ones ran. Re-running a migration can fail or corrupt data (0005 fails on a second run), and nothing stops a new version from starting against an outdated schema. Migrations 0005 and 0006 are the first ones that must reach an existing production database.

## What Changes

- The SQL migrations are embedded in the application binary.
- New `-migrate` flag: applies pending migrations, records the version in the database, and exits.
- New `-migrate-force <version>` flag: records a version without running anything, to baseline an existing database or recover from a failed migration.
- Docker Compose runs `-migrate` as a one-off service before the API starts.
- CI applies migrations with `-migrate` instead of `psql`.

## Capabilities

### New Capabilities
- `database-migrations`: How operators apply schema migrations and recover from failures.

### Modified Capabilities

None.

## Impact

- **Operations**: Existing databases need a one-time baseline (`-migrate-force 4`) before the first `-migrate`. Deployments add a `migrate` service to the Compose file.
- **Dependencies**: adds `github.com/golang-migrate/migrate/v4` with its pgx v5 driver.
- **Code**: new `migrations` embed package, `internal/shared/dbmigrate`, flags in `cmd/app/main.go`, CI workflow, README.
