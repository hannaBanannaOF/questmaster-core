# Tasks

## 1. Storage

- [ ] 1.1 Add migration `0007_recent_access` with `campaign_access` and `character_access` (primary key user and item, `accessed_at`, foreign keys `ON DELETE CASCADE`) and its down migration; verify with the dbmigrate tests
- [ ] 1.2 Add access repositories with an upsert that sets `accessed_at = now()`; verify with repository tests (CI Postgres) for first access, repeated access and cascade on campaign and character delete

## 2. Recording

- [ ] 2.1 Record the access in the campaign details use case after `CanView` passes, logging and ignoring a failed upsert; verify with use case tests for DM, player, non-member (nothing recorded) and a failing recorder (details still returned)
- [ ] 2.2 Record the access in the character details use case only when the requester is the character's player; verify with use case tests for the owner and for another user

## 3. Sort and field

- [ ] 3.1 Parse `sort` (`name` default, `recent`) for both lists, returning `ErrInvalidParam` for other values; verify with unit tests
- [ ] 3.2 Join the access tables in both list queries for `ORDER BY accessed_at DESC NULLS LAST, unaccent(lower(name)), id` with `sort=recent` and for `last_accessed_at` in both sorts; verify with repository tests (CI Postgres) for order, never opened items, another user's access and combined filters
- [ ] 3.3 Add `last_accessed_at` to both list responses and mappers, and the `sort` parameter to the Swagger annotations; verify with mapper tests and that CI regenerates `docs/`

## Workflow follow-up

- Frontend: dashboard asks for `sort=recent` and shows "opened …" from `last_accessed_at`.
- Archive the change after CI passes.
