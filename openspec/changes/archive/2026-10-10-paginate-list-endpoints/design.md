# Design

## Context

See proposal.md for motivation. The campaign list use case calls `GetByDmId` and `GetByPlayerId` and merges the results in Go, skipping duplicates. The character list supports the `game_system` and `without_campaign` filters through `AppContext.SetFilters`. Neither query has an `ORDER BY`.

## Goals / Non-Goals

**Goals:**
- Let clients load a page of each list, with the total, in a stable order.
- Let clients filter the campaign list by role and status.

**Non-Goals:**
- Cursor pagination or client-chosen sorting.
- A dashboard endpoint: clients compose the dashboard from the list endpoints.

## Decisions

**Offset pagination with an envelope.** `limit`/`offset` plus `{ items, total }`. Alternative: cursor pagination. Rejected for now because lists are per user and small enough that `OFFSET` stays cheap, and clients need `total` for counts and page numbers anyway. Alternative: keep the bare array and send the total in an `X-Total-Count` header. Rejected because headers are easy to lose in clients and proxies, and the frontend is the only client, so the breaking change is cheap now.

**One query for the campaign list.** `WHERE c.dm_id = $1 OR EXISTS (character of $1 in c)`, with the `role` filter choosing one side. Merging two queries in Go cannot apply `LIMIT`, `OFFSET` and `total` across both roles correctly.

**Order by name, then id.** `ORDER BY unaccent(lower(name)), id`, using the `unaccent` extension already installed by the init migration. Id breaks ties so pages never repeat or skip items.

**Defaults are safe, not unlimited.** Without `limit`, a request returns the first 50 items. Clients that need everything page through the list. A maximum of 100 keeps a single request bounded.

**Validation errors are 400.** `limit` outside 1..100, a negative `offset`, an unknown `role` or `status` return 400 through the existing `ErrInvalidParam`, like the current `game_system` filter.

## Risks / Trade-offs

- [The frontend breaks if released before or after the core] → Release both together; the proposal is marked BREAKING.
- [`total` costs a second count query] → Run it with the same filters; lists are per user, so the count stays small. Use `COUNT(*) OVER ()` if profiling shows a need.
