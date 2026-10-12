# Design

## Context

See proposal.md (Why). `character_sheet.player_id` is nullable, but `Character.PlayerID` is a non-pointer `userDomain.UserID` and the pg mapper builds it from the row unconditionally, so a row with a null player cannot be read today. Accepting an invite calls the character module's `LinkToCampaign(campaignID, characterSlug, userID)`, which loads the character, checks ownership, campaign and system, and sets `campaign_id`. The campaign details load the campaign's characters by `campaign_id`; the player count and membership checks filter by `player_id`, so null players never count or grant access.

## Goals / Non-Goals

**Goals:**
- One character table and one character type for both kinds; a premade character is simply a character with a campaign and no player.
- Claiming is atomic: two players can never both own the same premade character.

**Non-Goals:**
- Editing a character's name or max HP.
- Reserving a premade character for a given user.

## Decisions

**`PlayerID` becomes optional.** `Character.PlayerID` turns into `*userDomain.UserID`, the pg row uses `uuid.NullUUID`, and `IsPlayer` returns false when there is no player. `IsPremade()` is `PlayerID == nil`. Alternative: a separate `premade_character` table. Rejected because claiming would move rows between tables and lose the id and slug the party links use.

**Creation lives in the character module.** `CreatePremadeCharacter` takes the campaign through a port (as `UpdateHP` already does for `CampaignAccess`), checks `IsDM` and that the campaign is not ARCHIVED, and inserts with the campaign's system, `campaign_id` set and `player_id` null. The handler is mounted at `POST /campaign/{campaignID}/character` with the existing `CampaignIDMiddleware`, and the body has `name` and optional `max_hp`, validated like character creation.

**Claiming extends `LinkToCampaign`.** When the loaded character is premade, it must belong to the invite's campaign (else 400); the repository then runs `UPDATE character_sheet SET player_id = $1 WHERE id = $2 AND campaign_id = $3 AND player_id IS NULL` and zero affected rows means `ErrPremadeAlreadyClaimed` (409). The ownership, campaign and system checks for the requester's own characters stay as they are. Alternative: `SELECT ... FOR UPDATE` in a transaction. Rejected as more code for the same guarantee.

**Listing.** The invite details and the campaign details both read unclaimed premade characters through one repository query (`campaign_id = $1 AND player_id IS NULL ORDER BY name`). The campaign party query adds `player_id IS NOT NULL`. The user's character list already filters by `player_id = $1`, so premade characters never appear there.

**Deleting.** `CanDelete` also allows the campaign DM when the character is premade, so the delete use case loads the campaign when the character has one. Campaign deletion deletes `player_id IS NULL` characters of the campaign before unlinking the rest, in the same transaction.

## Risks / Trade-offs

- [Code that reads `PlayerID` directly breaks at compile time] → Intended: the compiler lists every place to update; the call sites are few (policy, mapper, read models).
- [Deleting the campaign deletes the DM's prepared characters] → They only make sense inside that campaign; claimed ones are kept like any player character.
