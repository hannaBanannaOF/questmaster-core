# database-migrations Specification

## Purpose
Defines how operators apply database schema migrations with the application binary, so every migration runs exactly once per database and failures are recoverable.

## Requirements

### Requirement: Migrate applies pending migrations once
Running the binary with `-migrate` SHALL apply, in order, every migration newer than the version recorded in the database, record the new version, and exit with status 0. When the database is already at the latest version it MUST change nothing and exit with status 0.

#### Scenario: Empty database
- **WHEN** an operator runs `-migrate` against an empty database
- **THEN** every migration is applied and the recorded version is the latest migration

#### Scenario: Database already up to date
- **WHEN** an operator runs `-migrate` against a database at the latest version
- **THEN** nothing is applied and the command exits with status 0

### Requirement: Force records a version without running migrations
Running the binary with `-migrate-force <version>` SHALL record that version as applied and clear any failed state, without running any migration, and exit with status 0.

#### Scenario: Baseline an existing database
- **WHEN** an operator runs `-migrate-force 4` against a database where migrations 0001 to 0004 were applied by hand, then runs `-migrate`
- **THEN** only the migrations after 0004 are applied

### Requirement: A failed migration stops the deploy
When a migration fails, `-migrate` SHALL exit with a non-zero status and leave the database marked as failed at that version, so the API service is not started and later `-migrate` runs refuse to continue until an operator forces a version.

#### Scenario: Migration error
- **WHEN** a migration fails while running `-migrate`
- **THEN** the command exits with a non-zero status

### Requirement: Every migration can be rolled back
Each migration SHALL have a down migration that reverts it, so that rolling back every migration and applying them again succeeds.

#### Scenario: Full rollback and reapply
- **WHEN** every migration is applied, then rolled back, then applied again
- **THEN** each step succeeds and the recorded version is the latest migration
