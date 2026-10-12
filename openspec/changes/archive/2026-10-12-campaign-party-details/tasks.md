# Tasks

## 1. Campaign characters

- [x] 1.1 Add `Slug`, `MaxHP` and `IsMine` to `CampaignCharacterReadModel`, filled in the campaign details use case from the loaded characters and the requester; verify with use case tests for the DM, a player with one character and a character without HP
- [x] 1.2 Add `slug`, `max_hp` and `is_mine` to `CampaignDetailResponseCharacterItem` and its mapper, and update the Swagger annotation; verify with a mapper test

## 2. Character access

- [x] 2.1 Add the `Character.CanView` policy (player; campaign DM and members when the character has a campaign) returning `ErrForbidden`; verify with domain tests for player, DM, member, non-member and character without a campaign
- [x] 2.2 Check `CanView` in the character details use case before recording the access, loading the campaign and the requester's membership through ports; verify with use case tests for each case, including that a forbidden read records no access
- [x] 2.3 Add the 403 response to the character details Swagger annotation and verify CI regenerates `docs/`

## Workflow follow-up

- Frontend: campaign page renders party cards linked to each sheet, with HP and "Sua ficha".
- Archive the change after CI passes.
