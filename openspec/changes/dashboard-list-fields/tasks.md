# Tasks

## 1. Character list

- [ ] 1.1 Join the linked campaign (slug, name) in the character list query and add `campaign` (nullable object) to the list item DTO and Swagger annotation; verify with a repository test for a linked and an unlinked character

## 2. Campaign list

- [ ] 2.1 Parse `without_players` (`true` filters, `false` is no filter, 400 otherwise) and add the player-count condition to the list query, combined with `role`, `status`, `sort` and `total`; verify with a handler test for the invalid value and a repository test for a DM with one campaign with players and one without

## Workflow follow-up

- Frontend: `dashboard-redesign` reads `campaign` and `without_players` (it already tolerates their absence).
- Archive the change after CI passes.
