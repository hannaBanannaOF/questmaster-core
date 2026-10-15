# Design

## Context

See proposal.md (Why). Campaign membership is derived from characters: a player is a member while one of their characters has the campaign's `campaign_id`. Accepting an invite links (or, for a premade character, claims) one character. Leaving or being removed unlinks it; deleting a campaign deletes its unclaimed premade characters and unlinks the rest. Characters without a campaign are only visible to their player.

## Goals / Non-Goals

**Goals:**
- Show the DM's and the players' names on the campaign page without any identity data in the core.
- Collect each name at the moment and for the purpose it is shown.

**Non-Goals:**
- A global nickname or profile; reading names from Keycloak.
- Names for characters outside a campaign: only their player sees them, so there is no one to show a name to.

## Decisions

**Stored on the records, not in a table of users.** `campaign.dm_name` belongs to the campaign and `character_sheet.player_name` to the character's stay in a campaign. The name is scoped to one table (the same person can be "Mestre Ana" in one campaign and "Ana" in another) and disappears with the record or the membership, so there is nothing to clean up separately. Alternative: one `user_display_name` table. Rejected: a single global name, and an extra deletion path for data that outlives every campaign.

**Required when created; edited into 0001.** The service has no real data yet, so the columns go straight into `0001_init` and the database is recreated, one last time. `campaign.dm_name` is `NOT NULL`, since every campaign is created with it; `character_sheet.player_name` is nullable, set only while the character is in a campaign.

**Cleared on unlink.** The repository updates that unlink a character (leave/remove, and the campaign deletion transaction) also set `player_name = NULL`, so a returning player names themselves again. The claim and link updates set it.

**Value object.** A `TableName` (trimmed, 1 to 50 runes) shared by both fields, validated in the transport layer like `CharacterName`.

**Editing.** `PATCH /campaign/{id}` is new and only accepts `dm_name` (DM only, `CanEdit`). `PATCH /character/{id}` accepts `name` and/or `player_name`: renaming keeps its rule (player or campaign DM), while `player_name` can only be changed by the character's player and only while the character is in a campaign (400 otherwise). A request with neither field returns 400.

**Suggestion comes from the client.** The frontend prefills the forms with the first name from the user's own session; the user confirms or edits it before it is sent. The core never reads another user's identity data.

## Risks / Trade-offs

- [A player with two characters in the same campaign can give two different names] → Each character shows its own; rare, and accepting it avoids a separate membership table.
- [Typing the name again at every table] → The session prefill makes it one confirmation; suggesting the last name used can come later.
- [Breaking change on create campaign and accept invite] → The frontend ships the new form fields together with this change.
