# character-details Specification

## Purpose
Defines who can read a character's details, so a character sheet only reaches its player and the members of its campaign.

## Requirements

### Requirement: Character details are visible to the player and campaign members
The system SHALL return a character's details to its player. When the character belongs to a campaign, the system SHALL also return them to the campaign's DM and to users with a character in the same campaign. Any other authenticated user MUST receive 403, and no access is recorded. An unknown character MUST return 404.

#### Scenario: Player reads their character
- **WHEN** a player reads the details of their own character
- **THEN** the system returns 200 with `is_player` true

#### Scenario: Campaign DM reads a player's character
- **WHEN** the DM of a character's campaign reads that character's details
- **THEN** the system returns 200 with `is_player` false

#### Scenario: Party member reads another player's character
- **WHEN** a user with a character in the same campaign reads another player's character
- **THEN** the system returns 200 with `is_player` false

#### Scenario: Non-member is forbidden
- **WHEN** a user who is neither the player, the campaign DM nor a campaign member reads a character
- **THEN** the system returns 403

#### Scenario: Character without a campaign
- **WHEN** a user other than the player reads a character that is not in any campaign
- **THEN** the system returns 403

#### Scenario: Unknown character
- **WHEN** a user reads a character id that does not exist
- **THEN** the system returns 404
