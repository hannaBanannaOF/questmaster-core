# Tasks

## 1. Table-name route

- [x] 1.1 Remove `PATCH /campaign/{id}` with its handler, request, response and Swagger annotations
- [x] 1.2 Replace `UpdateCampaignDmNameUseCase` with `UpdateTableNameUseCase` (DM → `dm_name`, player → their character's `player_name`, otherwise 403) and add the `CampaignPlayerNameUpdater` port with the repository `UpdatePlayerName(user, campaign, name)`; verify with use case tests for the DM, a player, an outsider (403) and an unknown campaign (404)
- [x] 1.3 Add `PATCH /campaign/{id}/table-name` (`{ "name" }`, validated with `NewTableName`) with Swagger annotations and wiring in module and bootstrap; verify with a route test for a valid name, an invalid name (400) and `PATCH /campaign/{id}` returning 404
- [x] 1.4 Verify `UpdatePlayerName` with a repository test changing only the target character

## 2. Character edit

- [x] 2.1 Make `name` required again in the character edit and drop `player_name` from command, request, response, read model and Swagger; delete `SetPlayerName`, `ErrNotInAnyCampaign`, `ErrNothingToEdit` and their error mappings; `UpdateDetails` no longer writes `player_name`; update the edit tests

## 3. Wrap-up

- [x] 3.1 Run `go vet` and the tests, validate the change with `openspec validate --strict`

## Workflow follow-up

- Frontend: move the table-name editing to `PATCH /campaign/{id}/table-name` and stop sending `player_name` to the character edit.
- Archive the change after CI passes.
