# Proposal

## Why

`player_count` is wrong. In the list of campaigns where the user plays, the query filters the joined characters to the requester's own before counting, so it only counts the requester's characters. Everywhere else it counts characters, not players, so a player with two characters counts twice.

## What Changes

- `player_count` is the number of distinct players with at least one character in the campaign, in every response that includes it (campaign list, invite details).
- Documents the existing rules of the current user's campaign list.

## Capabilities

### New Capabilities
- `campaign-list`: Which campaigns the current user sees in their list and what the player count means.

### Modified Capabilities

None.

## Impact

- **API**: `GET /core/api/v1/campaign` and `GET /core/api/v1/invite/{inviteHash}` return the corrected `player_count`. No contract change.
- **Code**: campaign PostgreSQL repository queries.
