# Design

## Context

Each campaign query joined `character_sheet` and grouped by campaign to compute `COUNT(cs.id)`. The player list query also filtered that join by `cs.player_id`, which reduced the count to the requester's characters.

## Goals / Non-Goals

**Goals:**
- One definition of `player_count`, shared by every campaign query.

**Non-Goals:**
- Storing the count (denormalization); the per-campaign subquery is cheap at the current scale.

## Decisions

**Correlated subquery in a shared `SELECT`.** All read queries use a `selectCampaign` constant whose `player_count` is `COUNT(DISTINCT player_id)` over the campaign's characters, so each query only adds its `WHERE`. The player list filters with `EXISTS` on the requester's characters instead of the join, so the filter can't affect the count. Alternative: keep the join with a second, unfiltered join for counting. Rejected as harder to read and easy to break again.

**`UpdateStatus` returns the real count.** Its `RETURNING` used a constant `0`; it now uses the same subquery, so the returned campaign is correct even though the handler only reads the status today.

## Risks / Trade-offs

- [The subquery runs once per listed campaign] → Acceptable for per-user lists; an index on `character_sheet.campaign_id` would help if lists grow, and is left for a later change.
