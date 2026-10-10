# Tasks

## 1. Query

- [ ] 1.1 Load the requester's characters for the campaigns of the current page (one query by campaign ids and player id, ordered by name) and attach them to the list read model; verify with repository tests (CI Postgres) for one character, several characters, and the DM of a campaign with players

## 2. Transport

- [ ] 2.1 Add `my_characters` to the campaign list response and mapper, and update the Swagger annotation; verify with a mapper test, including an empty list rendered as `[]`, not `null`

## Workflow follow-up

- Implement after `paginate-list-endpoints`, so the characters are loaded only for the returned page.
- Archive the change after CI passes.
