# Design

## Context

`dm_name` lives on `campaign` and `player_name` on `character_sheet` (see the archived `campaign-member-names` change). Today `UpdateCampaignDmNameUseCase` handles the DM path and `Character.SetPlayerName` the player path inside the character edit. The campaign use cases already reach characters through ports declared in `internal/campaign/app/usecases/ports.go`, served by character-module use cases.

## Goals / Non-Goals

**Goals:**
- One self-service route for table names; each route has a single permission rule.

**Non-Goals:**
- Editing campaign name or overview (`PATCH /campaign/{id}` comes back when that is needed).
- Any schema change.

## Decisions

- **Role decides the target, not the body.** The request carries only `name`; the use case loads the campaign (404), and if the requester is the DM it calls `UpdateDmName`, otherwise it updates the `player_name` of the requester's character in that campaign in one statement (`WHERE player_id = user AND campaign_id = campaign`); no row changed means the requester has no character there and returns `ErrNotCampaignMember` (403). One statement avoids a find-then-update race. Alternative: a `role` field in the body, rejected because the server already knows the role and the client could not choose anything else.
- **The DM-name use case becomes `UpdateTableNameUseCase`** in the campaign module, replacing `UpdateCampaignDmNameUseCase`, with a new consumer-declared port `CampaignPlayerNameUpdater.UpdatePlayerName(userID, campaignID, name) (bool, error)`, served by a character-module use case over the repository method of the same name. Alternative: reuse `UpdateDetails`, rejected because it also rewrites the name and slug.
- **Allowed in any status**, including archived: renaming yourself does not change play.
- **Character edit goes back to a required `name`**: the command and request use a plain `Name`; `SetPlayerName`, `ErrNotInAnyCampaign` and `ErrNothingToEdit` are deleted, and `UpdateDetails` stops writing `player_name`.
- **Response** `{ "name": "..." }` for the new route, mirroring what was saved.

## Risks / Trade-offs

- Breaking for the frontend, which already calls `PATCH /campaign/{id}` and sends `player_name` to the character edit → ship together with the frontend change.
- A DM is never a player in their own campaign (invites reject the DM), so the DM/player branch is unambiguous.
