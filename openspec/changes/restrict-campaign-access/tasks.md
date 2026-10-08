# Tasks

## 1. Campaign details access

- [x] 1.1 Add `ErrNotCampaignMember` and `Campaign.CanView`, map the error to 403, and verify the "outsider is forbidden" test in `get_campaign_details_test.go` passes
- [x] 1.2 Check membership in the details use case and load the invite only for the DM, and verify the DM, player and not-found tests in `get_campaign_details_test.go` pass
- [x] 1.3 Document the 403 and invite hash visibility in the Swagger annotations of `GetCampaignDetails` and verify CI regenerates `docs/`

## 2. Campaign invites

- [x] 2.1 Add `UserID` to `CreateInviteCommand`, check `CanEdit` in the create invite use case, wire the campaign finder in the invite module, and verify all `TestCreateInvite` cases pass
- [x] 2.2 Document the 403 and 404 in the Swagger annotations of `CreateInvite` and verify CI regenerates `docs/`

## 3. Character HP

- [x] 3.1 Add `ErrCharacterWithoutHP`, reject HP updates for characters without HP, map the error to 400, and verify all `TestUpdateHP` cases pass
- [x] 3.2 Return an error instead of panicking when a character row has only one HP column set, and verify `TestMapRowToDomainHP` passes

## 4. User profile

- [x] 4.1 Return null name and surname when the token has no given name, and verify `TestMapUserToUserResponse` passes

## 5. API errors

- [x] 5.1 Return a generic message for 500 responses and verify `TestFromHidesInternalErrors` passes

## 6. Integration

- [x] 6.1 Verify `go vet ./...` and `go test ./...` pass in CI

## Workflow follow-up

- Archive the change after CI passes, so the specs move to `openspec/specs/`.
