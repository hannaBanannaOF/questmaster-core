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
