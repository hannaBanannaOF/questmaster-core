# Design

## Context

See proposal.md for motivation. The campaign details use case already loads the campaign's characters (`FindByCampaignId`, selecting `cs.*`), decides membership with `Campaign.CanView(userID, hasCharacter)`, and maps them to `CampaignCharacterReadModel` with only id, name and current HP. The character details use case loads the character by id and returns it without checking who asks; it only records the access when the requester is the player. `Character.CanUpdate` already decides HP updates from the player and a `CampaignAccess` (the campaign's DM).

## Goals / Non-Goals

**Goals:**
- Return what the frontend needs to render party cards, without a request per character.
- Close the character details access gap with the same membership rule as the campaign details.

**Non-Goals:**
- A dedicated party endpoint.
- Exposing other players' identities.

## Decisions

**Extend the existing read model.** `CampaignCharacterReadModel` gains `Slug`, `MaxHP` and `IsMine`, filled from the characters the use case already loads, comparing `PlayerID` with the requester. Alternative: a `/campaign/{id}/characters` endpoint. Rejected because it adds a request and duplicates the membership check for data already in memory.

**Character access is campaign membership.** A new `Character.CanView(userID, campaign)` policy allows the player; and, when the character has a campaign, the campaign's DM and any user with a character in that campaign. The use case loads the campaign through the existing campaign finder port and asks the character repository whether the requester has a character in it. A forbidden read returns `ErrForbidden` (403), recorded nothing, like the campaign details.

**403, not 404, for non-members.** Same as the campaign details (specs/campaign-details), so clients can tell "not allowed" from "does not exist". The frontend shows its existing not-found page for both.

## Risks / Trade-offs

- [Clients that read arbitrary characters break with 403] → Only the frontend calls this endpoint, and it only links to characters the user may see (own list, campaign party). Marked BREAKING.
- [Two extra queries on character details] → Campaign and membership lookups by indexed ids; only for characters that have a campaign.
