# Tasks

## 1. Use case

- [x] 1.1 Add `IsDM` to `InviteDetailReadModel` and set it from `campaign.IsDM(cmd.UserID)` in the get invite detail use case
- [x] 1.2 Add `get_invite_detail_test.go` covering a player (`IsDM` false), the campaign DM (`IsDM` true) and an unknown invite (`ErrInviteNotFound`), and verify it passes

## 2. Transport

- [x] 2.1 Add `is_dm` to `InviteDetailsResponse` and its mapper, and verify the handler response with a mapper test
- [x] 2.2 Check the Swagger annotation of get invite details still describes the response, and verify CI regenerates `docs/`

## Workflow follow-up

- Archive the change after CI passes.
- Frontend: use `is_dm` on the invite page to show the DM a notice and a link to the campaign instead of the character picker.
