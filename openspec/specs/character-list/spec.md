# character-list Specification

## Purpose
Defines which characters appear in the current user's character list, how it is filtered, ordered and paginated.

## Requirements

### Requirement: The character list contains only the user's characters
The system SHALL list only the characters whose player is the requester. The `game_system` filter keeps characters of that game system, and `without_campaign=true` keeps characters not linked to any campaign. An unknown `game_system` MUST return 400.

#### Scenario: Own characters only
- **WHEN** a user requests their character list
- **THEN** every item is a character of that user, and no character of another player appears

#### Scenario: Characters without a campaign
- **WHEN** a user with one linked and one unlinked character requests their list with `without_campaign=true`
- **THEN** `items` contains only the unlinked character

#### Scenario: Unknown game system
- **WHEN** a user requests their character list with `game_system=UNKNOWN`
- **THEN** the system returns 400

### Requirement: The character list is ordered and paginated
The character list SHALL be ordered by character name, ignoring case and accents, then by id. It SHALL accept `limit` (1 to 100, default 50) and `offset` (0 or more, default 0), and respond with the page in `items` and the number of characters matching the filters in `total`. Invalid `limit` or `offset` values MUST return 400.

#### Scenario: Preview of the first characters
- **WHEN** a user with 7 characters requests their character list with `limit=5`
- **THEN** `items` contains the first 5 characters by name and `total` is 7

#### Scenario: Default page size
- **WHEN** a user with 60 characters requests their character list without `limit`
- **THEN** `items` contains 50 characters and `total` is 60

#### Scenario: Negative offset
- **WHEN** a user requests their character list with `offset=-1`
- **THEN** the system returns 400

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
