# campaign-invites Specification

## Purpose
Defines who can create a campaign invite and the possible outcomes, since an invite hash lets its holder join the campaign.

## Requirements

### Requirement: Only the campaign DM can create an invite
The system SHALL create a campaign invite only when the requester is the campaign DM. Any other authenticated user MUST receive 403 and no invite is created.

#### Scenario: DM creates an invite
- **WHEN** the campaign DM creates an invite for a campaign that has none
- **THEN** the system returns 201 with the invite hash

#### Scenario: Non-DM is forbidden
- **WHEN** a user who is not the campaign DM creates an invite for the campaign
- **THEN** the system returns 403
- **AND** no invite is created

#### Scenario: Campaign does not exist
- **WHEN** a user creates an invite for a campaign that does not exist
- **THEN** the system returns 404

### Requirement: A campaign has at most one invite
The system SHALL keep at most one invite per campaign. Creating an invite for a campaign that already has one MUST return 409.

#### Scenario: Invite already exists
- **WHEN** the campaign DM creates an invite for a campaign that already has an invite
- **THEN** the system returns 409

### Requirement: Archived campaigns do not take invites
The system SHALL reject creating or accepting an invite for an ARCHIVED campaign with 400. DRAFT, ACTIVE and PAUSED campaigns accept invites.

#### Scenario: Creating an invite for an archived campaign
- **WHEN** the campaign DM creates an invite for an ARCHIVED campaign
- **THEN** the system returns 400
- **AND** no invite is created

#### Scenario: Accepting an invite for an archived campaign
- **WHEN** a player accepts the invite of an ARCHIVED campaign
- **THEN** the system returns 400
- **AND** the character is not linked to the campaign

#### Scenario: Accepting an invite for a draft campaign
- **WHEN** a player accepts the invite of a DRAFT campaign with an eligible character
- **THEN** the system returns 204 and the character is linked to the campaign

### Requirement: The DM cannot join their own campaign as a player
The system SHALL reject the campaign DM accepting their own campaign's invite with 403.

#### Scenario: DM accepts their own invite
- **WHEN** the campaign DM accepts their campaign's invite with one of their characters
- **THEN** the system returns 403
- **AND** the character is not linked to the campaign

### Requirement: Accepting an invite links an eligible character
Accepting an invite SHALL link the given character to the invite's campaign only if the character belongs to the requester, is not linked to any campaign, and has the campaign's game system. Otherwise the system MUST return 400 and change nothing. An unknown invite MUST return 404.

#### Scenario: Eligible character
- **WHEN** a player accepts an invite with their own unlinked character of the campaign's game system
- **THEN** the system returns 204 and the character is linked to the campaign

#### Scenario: Character of another player
- **WHEN** a user accepts an invite with a character that belongs to another player
- **THEN** the system returns 400 and the character is not linked

#### Scenario: Character already in a campaign
- **WHEN** a player accepts an invite with a character that is already linked to a campaign
- **THEN** the system returns 400 and the character stays in its current campaign

#### Scenario: Different game system
- **WHEN** a player accepts an invite with a character of a different game system than the campaign
- **THEN** the system returns 400 and the character is not linked

#### Scenario: Unknown invite
- **WHEN** a player accepts an invite hash that does not exist
- **THEN** the system returns 404
