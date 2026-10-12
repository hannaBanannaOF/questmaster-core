# Design

## Context

See proposal.md (Why). Characters link to a campaign through `character_sheet.campaign_id`; campaign membership, the player count and the campaign list are all derived from characters with `player_id` set and that `campaign_id`. Campaign deletion already unlinks characters (`ON DELETE SET NULL`). The character module already loads a campaign through a port for `UpdateHP` (`CampaignAccess.IsDM`) and for premade characters.

## Goals / Non-Goals

**Goals:**
- One endpoint for leaving and for the DM removing, with the policy in the character domain.

**Non-Goals:**
- Notifications or a removal history.

## Decisions

**Route under the campaign.** `DELETE /campaign/{campaignID}/character/{characterID}` names both ids, so the use case can tell "not in this campaign" (400) from "not found" (404), and mirrors `POST /campaign/{campaignID}/character` for premade characters. Alternative: `POST /character/{id}/leave`. Rejected because the DM's removal is not a "leave".

**Policy on the character.** `Character.CanLeaveCampaign(user, campaign)` allows the player or `campaign.IsDM(user)`, and rejects a premade character and an ARCHIVED campaign. The campaign port gains the status the check needs (`IsArchived()`), as the premade creation already does.

**Unlink only.** The repository sets `campaign_id = NULL` with `WHERE id = $1 AND campaign_id = $2`; zero affected rows (removed meanwhile) is treated as 400. HP, slug and recent-access rows stay: `campaign_access` for the leaving player is harmless, since the campaign list is derived from characters.

## Risks / Trade-offs

- [A claimed premade character leaves with its player] → Consistent with "once claimed, it belongs to the player"; the DM can create another premade character.
- [The DM removes a player by mistake] → The player can rejoin through the invite while the campaign is not archived.
