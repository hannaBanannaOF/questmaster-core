# Proposal

## Why

Table names are edited in two unrelated places: `dm_name` through a generic `PATCH /campaign/{id}` and `player_name` through `PATCH /character/{id}`. Both are the same thing, "how I am called at this table", and mixing `player_name` into the sheet edit puts two permission rules in one request (the DM may rename the sheet but not the player), so a DM sending both fields gets a 403 for the whole request.

## What Changes

- Add `PATCH /campaign/{id}/table-name` with body `{ "name": "..." }`: the requester changes their own table name in the campaign. The DM changes `dm_name`; a player changes the `player_name` of their character in that campaign. Nobody else's name can be changed through it.
- **BREAKING** Remove `PATCH /campaign/{id}`; it only changed `dm_name` and there are no other editable campaign fields yet.
- **BREAKING** `PATCH /character/{id}` edits only the sheet: `player_name` is no longer accepted and `name` is required again.
- `PATCH /campaign/{id}/status` is unchanged: a status change is a rule-checked transition, not a field edit.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `table-names`: members change their own table name through one campaign-scoped endpoint, for DM and players alike.
- `character-edit`: the character edit no longer takes `player_name`; `name` is required.

## Impact

- Routes: `cmd/app/routes/campaign.go` (remove `PATCH /:campaignID`, add `PATCH /:campaignID/table-name`).
- Campaign module: the DM-name use case, handler, request and response become the table-name ones; a new port reaches the requester's character in the campaign.
- Character module: edit command, request, response, use case and Swagger lose `player_name`; `SetPlayerName`, `ErrNotInAnyCampaign` and `ErrNothingToEdit` go away.
- No schema change, so no migration.
- Frontend: rename calls move to the new endpoint.
