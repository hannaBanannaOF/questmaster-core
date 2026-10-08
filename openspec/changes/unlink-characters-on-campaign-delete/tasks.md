# Tasks

## 1. Database

- [ ] 1.1 Add migration 0005 setting `ON DELETE SET NULL` on the character foreign key and `ON DELETE CASCADE` on the invite foreign key, with a down migration restoring `NO ACTION`, and verify the CI migration up/down/up step passes
- [ ] 1.2 Add a campaign repository test that deletes a campaign with linked characters and an invite, and verify in CI that the characters remain unlinked with the same player and the invite is gone

## 2. Application

- [ ] 2.1 Remove the invite deletion from the delete campaign use case, drop the `CampaignInviteDeleter` port, and update module and bootstrap wiring; verify `go build ./...` passes
- [ ] 2.2 Remove `DeleteInviteUseCase`, its module getter and `InviteRepository.DeleteByCampaignID`; verify `go vet ./...` passes
- [ ] 2.3 Add delete campaign use case tests for DM deletes DRAFT and ARCHIVED, non-DM forbidden, ACTIVE rejected and campaign not found, and verify they pass in CI
- [ ] 2.4 List migration 0005 in the README and verify the list matches the files in `migrations/`

## Workflow follow-up

- Archive the change after CI passes.
