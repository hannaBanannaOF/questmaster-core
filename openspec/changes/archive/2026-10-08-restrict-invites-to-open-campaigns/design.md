# Design

## Context

Creating an invite already loads the campaign (to check `CanEdit`). Accepting an invite only loads the invite and then runs a single conditional `UPDATE` that links the character if it belongs to the requester, is unlinked and matches the campaign's game system.

## Goals / Non-Goals

**Goals:**
- Status and DM rules live in the campaign domain, next to the existing policies.

**Non-Goals:**
- Hiding or invalidating invites of archived campaigns in the invite details endpoint.
- Limiting how many characters one player links to the same campaign.

## Decisions

**New campaign policies `CanInvite` and `CanJoin`.** `CanInvite(user)` combines `CanEdit` with the status rule; `CanJoin(user)` checks the status rule and that the user is not the DM. Both return new domain errors (`ErrCampaignArchived` → 400, `ErrDMCannotJoin` → 403), mapped in `httperrors`.

**Accept loads the campaign through the existing port.** The accept use case receives the `InviteCampaignFinder` already used by create and details, loads the campaign after the invite, and calls `CanJoin` before linking. Alternative: add the status and DM conditions to the linking `UPDATE`. Rejected because the `UPDATE` can only report "not linked", losing the distinct 400/403 errors the spec requires.

## Risks / Trade-offs

- [The status can change between `CanJoin` and the `UPDATE`] → A campaign archived in that window may still get one player. Acceptable for the cost; the conditional `UPDATE` still guarantees the character rules.
