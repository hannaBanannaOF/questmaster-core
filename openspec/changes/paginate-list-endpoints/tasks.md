# Tasks

## 1. Pagination and filters

- [ ] 1.1 Add a shared page type (`limit`, `offset`, defaults 50 and 0, maximum 100) and parse it with the campaign `role` and `status` filters in `AppContext.SetFilters`, returning `ErrInvalidParam` for invalid values; verify with unit tests for defaults, bounds and unknown values

## 2. Campaign list

- [ ] 2.1 Replace `GetByDmId` and `GetByPlayerId` in the list use case with one repository query that applies role, status, `ORDER BY unaccent(lower(name)), id`, `LIMIT`/`OFFSET` and returns the total; verify with repository tests (CI Postgres) for order, pages, total, both role filters, status filter and another user's campaigns
- [ ] 2.2 Return `{ items, total }` from `GET /campaign` and update its Swagger annotation and parameters; verify with a mapper test

## 3. Character list

- [ ] 3.1 Add `ORDER BY unaccent(lower(name)), id`, `LIMIT`/`OFFSET` and the total to the character list query, keeping the `game_system` and `without_campaign` filters; verify with repository tests (CI Postgres) for order, pages, default page size, filters and another player's characters
- [ ] 3.2 Return `{ items, total }` from `GET /character` and update its Swagger annotation and parameters; verify with a mapper test

## Workflow follow-up

- Release together with the frontend change that reads `items` and `total`.
- Archive the change after CI passes.
