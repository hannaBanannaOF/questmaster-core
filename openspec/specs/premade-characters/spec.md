# premade-characters Specification

## Purpose
Defines premade characters: characters the campaign DM prepares with no player, who creates and deletes them, and where they appear, so players can join a campaign without bringing their own sheet.

## Requirements

### Requirement: Only the campaign DM creates premade characters
The system SHALL create a premade character in a campaign only when the requester is the campaign DM. The character has the given name and optional max HP (current HP starts at max HP), the campaign's game system, the campaign as its campaign, and no player. Any other authenticated user MUST receive 403 and no character is created. Creating one in an ARCHIVED campaign MUST return 400. An unknown campaign MUST return 404.

#### Scenario: DM creates a premade character
- **WHEN** the campaign DM creates a premade character named "Harvey" with max HP 12 in a DRAFT Call of Cthulhu campaign
- **THEN** the system returns 201 with the character slug
- **AND** the character has no player, belongs to the campaign, has the campaign's game system and current HP 12

#### Scenario: Non-DM is forbidden
- **WHEN** a user who is not the campaign DM creates a premade character in the campaign
- **THEN** the system returns 403
- **AND** no character is created

#### Scenario: Archived campaign
- **WHEN** the campaign DM creates a premade character in an ARCHIVED campaign
- **THEN** the system returns 400
- **AND** no character is created

#### Scenario: Unknown campaign
- **WHEN** a user creates a premade character in a campaign that does not exist
- **THEN** the system returns 404

### Requirement: The campaign DM may delete an unclaimed premade character
The system SHALL let the campaign DM delete a premade character of their campaign while it has no player. Once a player claims it, only that player may delete it, like any other character. Any other authenticated user MUST receive 403 and the character is kept.

#### Scenario: DM deletes an unclaimed premade character
- **WHEN** the campaign DM deletes a premade character of their campaign that has no player
- **THEN** the system returns 204 and the character no longer exists

#### Scenario: DM cannot delete a claimed character
- **WHEN** the campaign DM deletes a premade character that a player has claimed
- **THEN** the system returns 403
- **AND** the character still exists, owned by the player

#### Scenario: Other user is forbidden
- **WHEN** a user who is not the campaign DM deletes an unclaimed premade character
- **THEN** the system returns 403
- **AND** the character still exists

### Requirement: Premade characters are not in any user's character list
A character without a player SHALL NOT appear in any user's character list, nor in the `my_characters` of the campaign list.

#### Scenario: Character list of the DM
- **WHEN** the DM who created a premade character requests their character list
- **THEN** the premade character is not in the list
