# Tasks

## 1. Domain

- [ ] 1.1 Add `ErrCampaignArchived` and `ErrDMCannotJoin`, the `CanInvite` and `CanJoin` campaign policies, and their HTTP mapping (400 and 403); verify with domain tests for each status and for the DM

## 2. Use cases

- [ ] 2.1 Use `CanInvite` in the create invite use case and verify `TestCreateInvite` covers the archived campaign
- [ ] 2.2 Load the campaign and check `CanJoin` in the accept invite use case, wire the campaign finder in the invite module, and verify accept tests for archived, draft, DM and unknown invite pass
- [ ] 2.3 Update the Swagger annotations of create and accept invite with the new 400 and 403 responses, and verify CI regenerates `docs/`

## 3. Character linking

- [ ] 3.1 Add repository tests for linking an eligible character, another player's character, an already linked character and a different game system, and verify they pass in CI

## Workflow follow-up

- Archive the change after CI passes.
