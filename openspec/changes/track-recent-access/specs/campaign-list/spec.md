# Spec Delta

## ADDED Requirements

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
