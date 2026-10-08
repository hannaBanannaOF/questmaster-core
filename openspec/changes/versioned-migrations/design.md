# Design

## Context

Migrations are plain `NNNN_name.up.sql` / `.down.sql` files, which already match golang-migrate's naming. The production image is `scratch`, so it has no shell and no `psql`. Deployment uses Docker Compose.

## Goals / Non-Goals

**Goals:**
- Run migrations from the same image and binary that runs the API.
- A failed migration keeps the previous API version running.

**Non-Goals:**
- Rolling back from the binary. Down migrations are tested in CI; a production rollback uses the golang-migrate CLI against the same files.
- Refusing to start the API when the schema is behind.

## Decisions

**golang-migrate with embedded files.** Its file naming, `schema_migrations` version table and Postgres advisory lock (safe against two concurrent runs) fit as-is. Migrations are embedded with `embed.FS` and read through its `iofs` source, so the `scratch` image needs nothing extra. Alternatives: goose (would need file renames or annotations) and a CLI container (a second image to keep in sync with the code).

**pgx v5 driver, selected by rewriting the URL scheme.** golang-migrate picks the driver from the URL scheme, so `postgres://` and `postgresql://` in `DB_URL` are rewritten to `pgx5://`. This reuses the pgx version the API already depends on, instead of pulling in `lib/pq`.

**A one-off Compose service, not migrate-on-startup.** A `migrate` service runs the image with `-migrate`, and the API depends on it with `condition: service_completed_successfully`. If a migration fails, Compose doesn't start the new API container. Migrating on API startup would turn a migration error into a crash loop of the API.

**CI uses the same path.** CI runs `go run ./cmd/app -migrate` for the repository tests, and a separate test applies, rolls back and reapplies everything in a throwaway database, so the shared test database is never torn down while other packages use it.

## Risks / Trade-offs

- [Running `-migrate` on an existing database without the baseline re-runs 0001 to 0004] → 0002 to 0004 use `IF NOT EXISTS`/`OR REPLACE` and would mostly no-op, but the README makes the baseline step mandatory and explicit.
- [`-migrate-force` with a wrong version skips or repeats migrations] → Documented as a baseline/recovery tool only, with the exact commands.

## Migration Plan

1. Deploy the new image without starting the API yet.
2. Once, on the existing database: `docker compose run --rm migrate -migrate-force 4`.
3. `docker compose up -d`: the `migrate` service applies 0005 and 0006, then the API starts.
4. Rollback: run the golang-migrate CLI `down 2` with the files from the repository, then deploy the previous image.
