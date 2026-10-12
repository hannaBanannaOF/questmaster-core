# Tasks

## 1. Policy

- [x] 1.1 Add `Character.CanLeaveCampaign` (player or campaign DM; not premade; campaign not ARCHIVED) and the campaign port's `IsArchived`; verify with domain tests for player, DM, other member, premade and archived

## 2. Use case and endpoint

- [x] 2.1 Add the `LeaveCampaign` use case: 404 for unknown character or campaign, 400 when the character is not in that campaign, then the policy, then the conditional unlink; verify with use case tests for each scenario
- [x] 2.2 Add the repository unlink (`campaign_id = NULL WHERE id AND campaign_id`); verify with a repository test that membership, the campaign list and the player count drop for a player's last character
- [x] 2.3 Add `DELETE /campaign/{campaignID}/character/{characterID}` (204) with the Swagger annotation; verify with a handler test

## Workflow follow-up

- Frontend: "Sair da campanha" for the player and "Remover" on party cards for the DM.
- Archive the change after CI passes.
