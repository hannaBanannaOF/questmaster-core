# Tasks

## 1. Storage and domain

- [x] 1.1 Add `campaign.dm_name` (`NOT NULL`) and nullable `character_sheet.player_name` to `0001_init` (the database is recreated once more); verify with the migration tests
- [x] 1.2 Add the `TableName` value object (trimmed, 1 to 50 runes) with its error mapped to 400; verify with unit tests for trimming, blank, 50 and 51 characters

## 2. Setting the names

- [x] 2.1 Require `dm_name` on campaign creation and store it; verify with a handler test for missing and valid names and a repository test reading it back
- [x] 2.2 Require `player_name` when accepting an invite and store it in the link and claim updates; verify with use case tests for both paths and a repository test
- [x] 2.3 Clear `player_name` in the unlink update and in the campaign deletion transaction; verify with repository tests for leaving and for deleting a campaign

## 3. Editing

- [x] 3.1 Add `PATCH /campaign/{id}` changing `dm_name`, DM only; verify with use case tests for the DM, a player (403) and an invalid name
- [x] 3.2 Accept `player_name` on `PATCH /character/{id}` (player only, character in a campaign) alongside `name`; verify with use case tests for the player, the DM (403), a character without a campaign (400) and a request with neither field (400)

## 4. Campaign details

- [x] 4.1 Return `dm_name` and `characters[].player_name` in the campaign details read model, response and Swagger annotations; verify with a use case test and a mapper test including `null`

## Workflow follow-up

- Frontend: ask for the table name when creating a campaign and accepting an invite (prefilled with the session's first name), show `dm_name` and `player_name` on the campaign page, and offer editing.
- Archive the change after CI passes.
