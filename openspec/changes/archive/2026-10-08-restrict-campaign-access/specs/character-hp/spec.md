# Spec Delta

## Purpose

Defines who can update a character's current hit points and which values are accepted, including characters from game systems without hit points.

## ADDED Requirements

### Requirement: HP can be updated by the character's player or campaign DM
The system SHALL allow a character's current HP to be updated by the character's player, or by the DM of the campaign the character belongs to. Any other authenticated user MUST receive 403.

#### Scenario: Player updates HP
- **WHEN** the character's player updates the character's current HP to a valid value
- **THEN** the system returns 200 with the new current HP

#### Scenario: Campaign DM updates HP
- **WHEN** the DM of the character's campaign updates the character's current HP to a valid value
- **THEN** the system returns 200 with the new current HP

#### Scenario: Other user is forbidden
- **WHEN** a user who is neither the character's player nor its campaign DM updates the character's HP
- **THEN** the system returns 403

### Requirement: Current HP must stay within bounds
The system SHALL reject a current HP below 0 or above the character's max HP with 400.

#### Scenario: HP above max
- **WHEN** the character's player sets current HP above the character's max HP
- **THEN** the system returns 400

#### Scenario: Negative HP
- **WHEN** the character's player sets current HP below 0
- **THEN** the system returns 400

### Requirement: Characters without HP cannot have HP updated
The system SHALL reject an HP update for a character that has no HP with 400.

#### Scenario: Character has no HP
- **WHEN** the character's player updates the HP of a character created without HP
- **THEN** the system returns 400
