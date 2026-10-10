# campaign-list Specification

## Purpose
Defines which campaigns appear in the current user's campaign list and how the player count of a campaign is computed.

## Requirements

### Requirement: The list contains the user's campaigns as DM and as player
The system SHALL list every campaign where the requester is the DM or has at least one character, each campaign once.

#### Scenario: DM and player campaigns
- **WHEN** a user who is DM of one campaign and has a character in another requests their campaign list
- **THEN** the list contains both campaigns, each once

#### Scenario: Several characters in the same campaign
- **WHEN** a user with two characters in the same campaign requests their campaign list
- **THEN** that campaign appears once

### Requirement: Player count counts distinct players
The player count of a campaign SHALL be the number of distinct players with at least one character in the campaign, regardless of who requests it.

#### Scenario: Player with several characters
- **WHEN** a campaign has two characters of one player and one character of another player
- **THEN** its player count is 2 in the DM's list, in each player's list and in the invite details

#### Scenario: Campaign without characters
- **WHEN** a campaign has no characters
- **THEN** its player count is 0

### Requirement: The campaign list is ordered and paginated
The campaign list SHALL be ordered by campaign name, ignoring case and accents, then by id. It SHALL accept `limit` (1 to 100, default 50) and `offset` (0 or more, default 0), and respond with the page in `items` and the number of campaigns matching the filters in `total`. Invalid `limit` or `offset` values MUST return 400. Each user only ever sees their own campaigns: pagination and filters never include campaigns where the requester is neither the DM nor a player.

#### Scenario: First page
- **WHEN** a user with 3 campaigns named "Beta", "alfa" and "Ômega" requests their campaign list with `limit=2`
- **THEN** the system returns 200 with "alfa" and "Beta" in `items`
- **AND** `total` is 3

#### Scenario: Next page
- **WHEN** the same user requests their campaign list with `limit=2&offset=2`
- **THEN** `items` contains only "Ômega" and `total` is 3

#### Scenario: Invalid limit
- **WHEN** a user requests their campaign list with `limit=0` or `limit=101`
- **THEN** the system returns 400

#### Scenario: Another user's campaigns are never listed
- **WHEN** a user requests their campaign list with any filters or page
- **THEN** no campaign appears where the user is neither the DM nor has a character

### Requirement: The campaign list can be filtered by role and status
The campaign list SHALL accept `role` (`dm` for campaigns the requester runs, `player` for campaigns where the requester has a character and is not the DM) and `status` (DRAFT, ACTIVE, PAUSED or ARCHIVED). Filters combine, and `total` reflects them. Unknown values MUST return 400.

#### Scenario: Active campaigns as a player
- **WHEN** a user requests their campaign list with `role=player&status=ACTIVE`
- **THEN** `items` contains only ACTIVE campaigns where the user has a character and is not the DM

#### Scenario: Campaigns as DM
- **WHEN** a user requests their campaign list with `role=dm`
- **THEN** `items` contains only campaigns the user is the DM of

#### Scenario: Unknown role
- **WHEN** a user requests their campaign list with `role=admin`
- **THEN** the system returns 400

### Requirement: Campaign counts by status
The system SHALL return, for any authenticated user, the number of their own campaigns in each status, with every status present (0 when there is none). The optional `role` filter (`dm` or `player`) restricts the counts like it restricts the campaign list. Campaigns where the requester is neither the DM nor a player MUST NOT be counted. An unknown `role` MUST return 400.

#### Scenario: DM counts
- **WHEN** a user who runs 2 ACTIVE campaigns and 1 ARCHIVED campaign requests the counts with `role=dm`
- **THEN** the system returns 200 with ACTIVE 2, DRAFT 0, PAUSED 0 and ARCHIVED 1

#### Scenario: Both roles
- **WHEN** a user who runs one DRAFT campaign and plays in one ACTIVE campaign requests the counts without `role`
- **THEN** the system returns DRAFT 1 and ACTIVE 1

#### Scenario: Another user's campaigns are not counted
- **WHEN** a user requests the counts
- **THEN** campaigns where the user is neither the DM nor has a character are not counted

#### Scenario: Unknown role
- **WHEN** a user requests the counts with `role=admin`
- **THEN** the system returns 400

### Requirement: The campaign list shows the requester's characters
Each campaign in the list SHALL include `my_characters`: the requester's characters linked to that campaign, with slug and name, ordered by name. Characters of other players MUST NOT appear in it, even for the campaign DM. It SHALL be empty when the requester has no character in the campaign.

#### Scenario: Player with one character
- **WHEN** a player with one character in a campaign requests their campaign list
- **THEN** that campaign's `my_characters` contains that character's slug and name

#### Scenario: Player with several characters
- **WHEN** a player with two characters in the same campaign requests their campaign list
- **THEN** that campaign appears once, with both characters in `my_characters`

#### Scenario: DM does not see players' characters
- **WHEN** the DM of a campaign with two players requests their campaign list
- **THEN** that campaign's `my_characters` is empty
