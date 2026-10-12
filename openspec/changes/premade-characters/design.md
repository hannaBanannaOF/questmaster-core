# Design

## Context

See proposal.md (Why). `character_sheet.player_id` is nullable, but `Character.PlayerID` is a non-pointer `userDomain.UserID` and the pg mapper builds it from the row unconditionally, so a row with a null player cannot be read today. Accepting an invite calls the character module's `LinkToCampaign(campaignID, characterSlug, userID)`, which runs a single conditional `UPDATE` (`UpdateCampaign`: slug, requester as player, no campaign, same system) without loading the character, and returns 400 when no row changes. Deleting a campaign is a single `DELETE` that relies on the foreign keys of migration 0005; no repository uses transactions today. The campaign details load the campaign's characters by `campaign_id`; the player count and membership checks filter by `player_id`, so null players never count or grant access.

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

**Claiming extends `LinkToCampaign`.** It first loads the character by slug (`FindBySlug`). A character with a player keeps the current conditional `UPDATE`, unchanged. A premade character must belong to the invite's campaign (else 400), and the requester must have no character in the campaign yet (checked first through `GetByPlayerInCampaigns`, else `ErrAlreadyInCampaign`, 409). The repository then runs `UPDATE character_sheet SET player_id = $1 WHERE id = $2 AND campaign_id = $3 AND player_id IS NULL AND NOT EXISTS (SELECT 1 FROM character_sheet o WHERE o.campaign_id = $3 AND o.player_id = $1)`, so concurrent claims cannot break either rule; zero affected rows means `ErrPremadeAlreadyClaimed` (409). Alternative: `SELECT ... FOR UPDATE` in a transaction. Rejected as more code for the same guarantee.

**Listing.** The invite details and the campaign details both read unclaimed premade characters through one repository query (`campaign_id = $1 AND player_id IS NULL ORDER BY name`). The campaign party query adds `player_id IS NOT NULL`. The user's character list already filters by `player_id = $1`, so premade characters never appear there.

**Deleting.** `CanDelete` also allows the campaign DM when the character is premade, so the delete use case loads the campaign when the character has one. `CampaignRepositoryPG.DeleteById` runs in a pgx transaction: it deletes the campaign's characters with `player_id IS NULL`, then the campaign, whose foreign keys (migration 0005) unlink the remaining characters and remove the invite.

## Risks / Trade-offs

- [Code that reads `PlayerID` directly breaks at compile time] → Intended: the compiler lists every place to update; the call sites are few (policy, mapper, read models).
- [Deleting the campaign deletes the DM's prepared characters] → They only make sense inside that campaign; claimed ones are kept like any player character.
