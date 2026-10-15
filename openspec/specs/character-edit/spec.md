# character-edit Specification

## Purpose
Defines who can rename a character after creation, and how its slug follows the new name.

## Requirements

### Requirement: Renaming regenerates the slug
Changing a character's name SHALL regenerate its slug from the new name, with the same rules and uniqueness as on creation. The old slug MUST no longer resolve to the character. Recent access and campaign links, which use the character id, are unaffected.

#### Scenario: Slug follows the new name
- **WHEN** a player renames the character "Maria", slug "maria", to "Jão"
- **THEN** the response slug is "jao" (or "jao-N" when taken)
- **AND** resolving "maria" returns 404

#### Scenario: Same name keeps the slug
- **WHEN** a player saves the character with its current name
- **THEN** its slug does not change

### Requirement: The player or the campaign DM renames a character
The system SHALL update a character's `name` when the requester is the character's player or the DM of the character's campaign, and return 200 with the character's slug and name. Any other authenticated user MUST receive 403 and nothing changes. An unknown character MUST return 404.

#### Scenario: Player renames their character
- **WHEN** a player changes their character's name to "Harvey Walters"
- **THEN** the system returns 200 with the new name

#### Scenario: DM edits an unclaimed premade character
- **WHEN** the campaign DM renames an unclaimed premade character of their campaign
- **THEN** the system returns 200 with the new name

#### Scenario: Other user is forbidden
- **WHEN** a user who is neither the character's player nor its campaign DM edits the character
- **THEN** the system returns 403
- **AND** the character is unchanged

#### Scenario: Unknown character
- **WHEN** a user edits a character id that does not exist
- **THEN** the system returns 404

### Requirement: The new name is validated
The name SHALL follow the same rules as character creation, and `player_name` the table name rules (see `table-names`). A request MUST include `name`, `player_name` or both; a request with neither, or with an invalid value, MUST return 400 and change nothing.

#### Scenario: Empty name
- **WHEN** the player sets the character's name to an empty string
- **THEN** the system returns 400 and the name is unchanged

#### Scenario: No name
- **WHEN** the player sends a request with neither `name` nor `player_name`
- **THEN** the system returns 400
