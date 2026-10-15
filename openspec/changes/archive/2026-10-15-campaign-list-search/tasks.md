# Tasks

## 1. Search

- [x] 1.1 Parse `q` in `filters.go` (trim, empty = none, more than 100 characters = 400); verify with a handler test for the blank and too-long cases
- [x] 1.2 Add the unaccented, case-insensitive "name contains" condition to the list and count queries, combined with the other filters; verify with a repository test for accents, case, a combined filter and another user's campaign

## Workflow follow-up

- Frontend: `campaign-list-redesign` already sends `q`.
- Archive the change after CI passes.
