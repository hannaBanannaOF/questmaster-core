# character-edit Specification

## Purpose
Defines who can edit a character's name and max HP after creation, and how the current HP and the slug follow the edit.

## Requirements

### Requirement: The player or the campaign DM edits a character's name and max HP
The system SHALL update a character's `name` and/or `max_hp` when the requester is the character's player or the DM of the character's campaign, and return 200 with the character's slug, name, current HP and max HP. Any other authenticated user MUST receive 403 and nothing changes. An unknown character MUST return 404.

#### Scenario: Player renames their character
- **WHEN** a player changes their character's name to "Harvey Walters"
- **THEN** the system returns 200 with the new name

#### Scenario: Campaign DM raises max HP
- **WHEN** the DM of a character's campaign sets its max HP from 10 to 13
- **THEN** the system returns 200 with max HP 13 and the current HP unchanged

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

### Requirement: Edited values are validated
The name SHALL follow the same rules as character creation, and max HP MUST be a positive integer. A request with neither field, or with an invalid value, MUST return 400 and change nothing.

#### Scenario: Empty name
- **WHEN** the player sets the character's name to an empty string
- **THEN** the system returns 400 and the name is unchanged

#### Scenario: Non-positive max HP
- **WHEN** the player sets max HP to 0
- **THEN** the system returns 400 and the HP is unchanged

#### Scenario: Nothing to change
- **WHEN** the player sends a request with neither `name` nor `max_hp`
- **THEN** the system returns 400

### Requirement: Current HP follows the max HP
When max HP is set below the current HP, the system SHALL lower the current HP to the new max. When max HP is set on a character without HP, the system SHALL give it HP with current HP equal to the max.

#### Scenario: Max HP below current HP
- **WHEN** the player sets max HP to 8 on a character with 10 / 12 HP
- **THEN** the character has 8 / 8 HP

#### Scenario: Character without HP gets HP
- **WHEN** the player sets max HP to 12 on a character created without HP
- **THEN** the character has 12 / 12 HP

### Requirement: Renaming regenerates the slug
Changing a character's name SHALL regenerate its slug from the new name, with the same rules and uniqueness as on creation. The old slug MUST no longer resolve to the character. Recent access and campaign links, which use the character id, are unaffected.

#### Scenario: Slug follows the new name
- **WHEN** a player renames the character "Maria", slug "maria", to "Jão"
- **THEN** the response slug is "jao" (or "jao-N" when taken)
- **AND** resolving "maria" returns 404

#### Scenario: Same name keeps the slug
- **WHEN** a player saves the character with its current name
- **THEN** its slug does not change
