# Spec Delta

## ADDED Requirements

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
