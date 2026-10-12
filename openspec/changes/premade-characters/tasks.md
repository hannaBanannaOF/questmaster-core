# Tasks

## 1. Characters without a player

- [ ] 1.1 Make `Character.PlayerID` optional (`*userDomain.UserID`), add `IsPremade()`, read `player_id` as `uuid.NullUUID` in the pg rows and mapper, and fix every caller; verify with a domain test for `IsPlayer`/`IsPremade` and a repository test reading a character without a player

## 2. Creating premade characters

- [ ] 2.1 Add the `CreatePremadeCharacter` use case (campaign DM only, not ARCHIVED, campaign's system, current HP = max HP) and its repository insert; verify with use case tests for DM, non-DM, archived and unknown campaign
- [ ] 2.2 Add `POST /campaign/{campaignID}/character` with request validation, 201 with the slug, and the Swagger annotation; verify with a handler test

## 3. Claiming through the invite

- [ ] 3.1 Extend `LinkToCampaign` to load the character by slug and claim an unclaimed premade character of the invite's campaign with a conditional update, returning `ErrAlreadyInCampaign` (409) when the requester already has a character in the campaign, `ErrPremadeAlreadyClaimed` (409) when no row changes and 400 for a premade character of another campaign; verify with use case tests (including a player who already has a character) and a repository test for two claims of the same character
- [ ] 3.2 Add `premade_characters` (slug, name, current and max HP, by name) to the invite details; verify with a use case test that claimed characters are left out

## 4. Campaign details, deletion and character list

- [ ] 4.1 Limit the campaign details `characters` to characters with a player and add `premade_characters` for the DM only (empty for players); verify with use case tests for the DM and a player, and a repository test that premade characters do not change the player count
- [ ] 4.2 Let the campaign DM delete an unclaimed premade character (`CanDelete`) and keep 403 for a claimed one and for other users; verify with domain and use case tests
- [ ] 4.3 Run `CampaignRepositoryPG.DeleteById` in a transaction that deletes unclaimed premade characters before the campaign; verify with a repository test for one claimed and one unclaimed character
- [ ] 4.4 Verify with a repository test that a premade character is not in the DM's character list nor in `my_characters`

## Workflow follow-up

- Frontend: DM creates and deletes premade characters on the campaign page; the invite page offers them.
- Archive the change after CI passes.
