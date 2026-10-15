# Spec Delta

## ADDED Requirements

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
The name SHALL follow the same rules as character creation. A request without a valid name MUST return 400 and change nothing.

#### Scenario: Empty name
- **WHEN** the player sets the character's name to an empty string
- **THEN** the system returns 400 and the name is unchanged

#### Scenario: No name
- **WHEN** the player sends a request without `name`
- **THEN** the system returns 400

## REMOVED Requirements

### Requirement: The player or the campaign DM edits a character's name and max HP
**Reason**: The core no longer stores HP, so editing only renames.
**Migration**: Replaced by "The player or the campaign DM renames a character"; max HP is edited through the game system's service.

### Requirement: Edited values are validated
**Reason**: Max HP is no longer an editable value.
**Migration**: Replaced by "The new name is validated".

### Requirement: Current HP follows the max HP
**Reason**: The core no longer stores HP.
**Migration**: The game system's service manages its characters' attributes.
