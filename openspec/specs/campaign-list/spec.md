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

### Requirement: Opening a campaign is recorded for its member
When the campaign DM or a player with a character in the campaign successfully reads the campaign details, the system SHALL record that moment as the requester's last access to that campaign, replacing the previous one. A read that returns 403 or 404 MUST NOT record anything. A failure to record MUST NOT fail the read.

#### Scenario: Player opens a campaign
- **WHEN** a player with a character in the campaign reads its details
- **THEN** the system returns 200 and the player's last access to that campaign is now

#### Scenario: Non-member is not recorded
- **WHEN** a user who is neither the DM nor has a character in the campaign reads its details
- **THEN** the system returns 403 and no access is recorded for that user

### Requirement: The campaign list can be sorted by recent access
The campaign list SHALL accept `sort` with `name` (default) or `recent`. With `recent`, campaigns SHALL be ordered by the requester's last access, most recent first; campaigns the requester never opened come after, ordered by name ignoring case and accents, then by id. Each item SHALL include `last_accessed_at`, the requester's last access or `null`, for both sorts. An unknown `sort` MUST return 400. Sorting never adds campaigns where the requester is neither the DM nor a player.

#### Scenario: Most recent first
- **WHEN** a user opened "Beta" yesterday and "alfa" today, never opened "Ômega", and requests their campaign list with `sort=recent`
- **THEN** `items` is "alfa", "Beta", "Ômega"
- **AND** "Ômega" has `last_accessed_at` null

#### Scenario: Recent access of another user does not count
- **WHEN** another member opened "Ômega" today and the requester never did
- **THEN** with `sort=recent`, "Ômega" stays after the campaigns the requester opened

#### Scenario: Unknown sort
- **WHEN** a user requests their campaign list with `sort=popular`
- **THEN** the system returns 400

### Requirement: The campaign list can be searched by name

The campaign list SHALL accept `q` and keep only campaigns whose name contains it, ignoring case and accents. Leading and trailing spaces SHALL be trimmed, and an empty `q` SHALL apply no filter. `q` combines with every other filter, sort and page, and `total` reflects it. A `q` longer than 100 characters MUST return 400. Searching never adds campaigns where the requester is neither the DM nor a player.

#### Scenario: Search ignores case and accents

- **WHEN** a user with campaigns "Arkham, 1925" and "Além das Montanhas" requests their campaign list with `q=ALEM`
- **THEN** `items` contains only "Além das Montanhas" and `total` is 1

#### Scenario: Search combines with filters

- **WHEN** a DM requests their campaign list with `role=dm&status=PAUSED&q=boston`
- **THEN** `items` contains only PAUSED campaigns they run whose name contains "boston"

#### Scenario: Other users' campaigns are not found

- **WHEN** a user searches for the name of a campaign where they are neither the DM nor a player
- **THEN** `items` is empty and `total` is 0

#### Scenario: Blank search

- **WHEN** a user requests their campaign list with `q=%20%20`
- **THEN** the result is the same as without `q`

#### Scenario: Search too long

- **WHEN** a user requests their campaign list with a `q` of 101 characters
- **THEN** the system returns 400
