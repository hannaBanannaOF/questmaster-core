# Tasks

## 1. Repository

- [x] 1.1 Replace the campaign read queries with a shared select that counts distinct players, filter the player list with `EXISTS`, and return the real count from `UpdateStatus`; verify `go vet ./...` passes
- [x] 1.2 Add repository tests for the player count in the DM list, each player's list and `FindById`, and for a campaign without characters, and verify they pass in CI

## 2. Use case

- [x] 2.1 Add a current user campaigns use case test for DM and player campaigns being listed once, and verify it passes in CI

## Workflow follow-up

- Archive the change after CI passes.
