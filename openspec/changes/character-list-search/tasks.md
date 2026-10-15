# Tasks

## 1. Search

- [ ] 1.1 Parse `q` on the character list (trim, empty = none, more than 100 characters = 400); verify with a handler test for the blank and too-long cases
- [ ] 1.2 Add the unaccented, case-insensitive "name contains" condition to the character list and count queries, reusing the normalization from `campaign-list-search`; verify with a repository test for accents, case, `without_campaign` combined and another player's character

## Workflow follow-up

- Frontend: `character-list-redesign` already sends `q`.
- Archive the change after CI passes.
