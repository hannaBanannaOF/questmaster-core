# Spec Delta

## ADDED Requirements

### Requirement: Opening an own character is recorded
When a user successfully reads the details of a character they play, the system SHALL record that moment as their last access to it, replacing the previous one. Reading another player's character, or a read that returns 403 or 404, MUST NOT record anything. A failure to record MUST NOT fail the read.

#### Scenario: Player opens their character
- **WHEN** a player reads the details of their own character
- **THEN** the system returns 200 and their last access to it is now

#### Scenario: Another player's character is not recorded
- **WHEN** a campaign DM reads the details of a player's character
- **THEN** no access is recorded for the DM

### Requirement: The character list can be sorted by recent access
The character list SHALL accept `sort` with `name` (default) or `recent`. With `recent`, characters SHALL be ordered by the requester's last access, most recent first; characters never opened come after, ordered by name ignoring case and accents, then by id. Each item SHALL include `last_accessed_at`, the requester's last access or `null`, for both sorts. The `game_system` and `without_campaign` filters combine with the sort. An unknown `sort` MUST return 400.

#### Scenario: Recent characters preview
- **WHEN** a user with 7 characters, 3 of them opened, requests their character list with `sort=recent&limit=5`
- **THEN** `items` starts with the 3 opened characters, most recent first, followed by 2 never opened by name
- **AND** `total` is 7

#### Scenario: Unknown sort
- **WHEN** a user requests their character list with `sort=popular`
- **THEN** the system returns 400
