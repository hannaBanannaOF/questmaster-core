# Design

## Context

See proposal.md for motivation. Authorization rules already live in the domain as policy methods (`Campaign.CanEdit`, `Campaign.CanDelete`, `Character.CanUpdate`), and use cases call them before persisting. The campaign domain cannot import the character domain, because the character domain already imports the campaign domain.

## Goals / Non-Goals

**Goals:**
- Enforce the access rules in specs/campaign-details and specs/campaign-invites in the domain and use cases, not in HTTP handlers.
- Turn the three nil pointer panics into defined responses.

**Non-Goals:**
- Checking campaign status when creating or accepting invites.
- Refactoring the shared error mapping so it no longer imports every module.

## Decisions

**Membership check takes a precomputed flag.** `Campaign.CanView(userID, hasCharacter bool)` receives whether the user has a character in the campaign, and the details use case computes it from the characters it already loads. Alternative: pass the characters, or an interface over them, into the campaign domain. Rejected because it creates an import cycle (or a new port only used once) for a single boolean.

**The invite is only loaded for the DM.** The details use case skips the invite lookup for non-DMs, instead of loading it and clearing the hash in the mapper. The hash never reaches the read model for non-DMs, and players save one query.

**Invite creation reuses `CanEdit`.** Creating an invite is a DM-only edit of the campaign, so the use case gets the campaign through the existing `InviteCampaignFinder` port and calls `CanEdit`, which returns `ErrNotDM` (403).

**HP rows with only one column set are an error, not a panic.** The row mapper returns an error when exactly one of `current_hp`/`max_hp` is null. The application never writes that state, so it signals corrupt data and correctly surfaces as a logged 500.

**500 responses use a fixed message.** `httperrors.From` returns "Internal server error" for unknown errors. `ErrorHandlerMiddleware` already logs the original error, so nothing is lost for debugging.

## Risks / Trade-offs

- [Clients that showed campaign details to non-members break with 403] → Non-members must use the invite details endpoint, which is unchanged. Called out as BREAKING in the proposal.
- [Players lose the invite hash and can no longer share invites] → Intended: sharing an invite is a DM decision.
- [A player who removes their last character loses access to the campaign details] → Matches the membership definition in specs/campaign-details. There is no detach endpoint today.
