# Spec Delta

## MODIFIED Requirements

### Requirement: The campaign list can be filtered by role and status
The campaign list SHALL accept `role` (`dm` for campaigns the requester runs, `player` for campaigns where the requester has a character and is not the DM), `status` (DRAFT, ACTIVE, PAUSED or ARCHIVED) and `without_players` (`true` keeps only campaigns whose player count is 0; `false` or absent applies no filter). Filters combine, and `total` reflects them. Unknown values MUST return 400.

#### Scenario: Active campaigns as a player
- **WHEN** a user requests their campaign list with `role=player&status=ACTIVE`
- **THEN** `items` contains only ACTIVE campaigns where the user has a character and is not the DM

#### Scenario: Campaigns as DM
- **WHEN** a user requests their campaign list with `role=dm`
- **THEN** `items` contains only campaigns the user is the DM of

#### Scenario: DM campaigns nobody joined
- **WHEN** a DM of two campaigns, one with 3 players and one with none, requests their campaign list with `role=dm&without_players=true`
- **THEN** `items` contains only the campaign with no players and `total` is 1

#### Scenario: Other users' campaigns without players are not listed
- **WHEN** a user requests their campaign list with `without_players=true`
- **THEN** campaigns where the user is neither the DM nor a player do not appear

#### Scenario: Unknown role
- **WHEN** a user requests their campaign list with `role=admin`
- **THEN** the system returns 400

#### Scenario: Invalid without_players
- **WHEN** a user requests their campaign list with `without_players=maybe`
- **THEN** the system returns 400
