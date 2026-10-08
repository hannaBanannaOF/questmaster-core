# Tasks

## 1. Database

- [x] 1.1 Add migration 0006 replacing both slug trigger functions with an empty-slug fallback and fixing existing empty slugs, with a down migration restoring the original functions; verify the CI migration up/down/up step passes
- [x] 1.2 Add repository tests for slug derivation, duplicate names and the fallback for campaigns and characters, and verify they pass in CI
- [x] 1.3 List migration 0006 in the README and verify the list matches the files in `migrations/`

## Workflow follow-up

- Archive the change after CI passes.
