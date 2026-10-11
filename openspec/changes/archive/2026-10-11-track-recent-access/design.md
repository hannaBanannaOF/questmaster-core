# Design

## Context

See proposal.md for motivation. Campaign details are readable by the DM and by players with a character in the campaign (specs/campaign-details); the use case already decides membership before returning. Character details can be read by more users than the owner, but the character list only ever contains the requester's own characters (specs/character-list). Both lists are single queries with `ORDER BY unaccent(lower(name)), id` and offset pagination.

## Goals / Non-Goals

**Goals:**
- Order the user's lists by when they last opened each item, across devices.
- Expose the last access time so clients can show it.

**Non-Goals:**
- Counting accesses or ranking by frequency.
- An activity feed or audit log.

## Decisions

**One table per item type, one row per user and item.** `campaign_access(user_id, campaign_id, accessed_at)` and `character_access(user_id, character_id, accessed_at)`, primary key `(user_id, <item>_id)`, foreign keys `ON DELETE CASCADE`. Each open is an `INSERT … ON CONFLICT DO UPDATE SET accessed_at = now()`. Alternative: one polymorphic table with an item type column. Rejected because it cannot have foreign keys, so deleted campaigns and characters would leave orphan rows.

**Recorded in the detail use cases, after authorization.** The campaign details use case records the access only after `CanView` passes; the character details use case records it only when the requester is the character's player, since only own characters appear in the character list. A 403 or 404 records nothing.

**Recording never fails the read.** If the upsert fails, the use case logs the error and still returns the details. Losing one recency update is harmless; failing to show a campaign is not.

**Recency sort joins the access table.** `sort=recent` becomes `LEFT JOIN <item>_access a ON a.user_id = $requester AND a.<item>_id = item.id ORDER BY a.accessed_at DESC NULLS LAST, unaccent(lower(name)), id`. The same join fills `last_accessed_at` for both sorts. The name and id tie-breakers keep pages stable.

**Validation is 400.** An unknown `sort` returns 400 through the existing `ErrInvalidParam`, like the other list filters.

## Risks / Trade-offs

- [A write on every detail read] → One indexed upsert per open; the detail read already runs several queries. Revisit with a short throttle (skip if updated in the last minute) only if profiling shows it.
- [Reads through the frontend's prefetching could count as opens] → The frontend loads details only on navigation, not on link hover; documented for clients.
