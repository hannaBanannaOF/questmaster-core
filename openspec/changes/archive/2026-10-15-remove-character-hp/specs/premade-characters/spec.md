# Spec Delta

## MODIFIED Requirements

### Requirement: Only the campaign DM creates premade characters
The system SHALL create a premade character in a campaign only when the requester is the campaign DM. The character has the given name, the campaign's game system, the campaign as its campaign, and no player. Any other authenticated user MUST receive 403 and no character is created. Creating one in an ARCHIVED campaign MUST return 400. An unknown campaign MUST return 404.

#### Scenario: DM creates a premade character
- **WHEN** the campaign DM creates a premade character named "Harvey" in a DRAFT Call of Cthulhu campaign
- **THEN** the system returns 201 with the character slug
- **AND** the character has no player, belongs to the campaign and has the campaign's game system

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
