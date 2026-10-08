# Tasks

## 1. Migration runner

- [x] 1.1 Add golang-migrate with the pgx v5 driver, embed `migrations/*.sql`, and add `internal/shared/dbmigrate` with `Up` and `Force`; verify `go build ./...` passes
- [x] 1.2 Add the `-migrate` and `-migrate-force` flags to `cmd/app/main.go`, running before the OIDC setup and exiting with a non-zero status on failure; verify `go vet ./...` passes
- [x] 1.3 Add a `dbmigrate` test that, in a throwaway database, applies everything, re-runs as a no-op, rolls everything back and reapplies, and baselines a hand-migrated database with force; verify it passes in CI

## 2. CI and documentation

- [x] 2.1 Replace the `psql` loop in CI with `go run ./cmd/app -migrate`, and verify the repository tests pass in CI
- [x] 2.2 Document the Compose `migrate` service, the one-time baseline, and failed-migration recovery in the README, and verify the documented flags match `cmd/app/main.go`

## Workflow follow-up

- Archive the change after CI passes.
