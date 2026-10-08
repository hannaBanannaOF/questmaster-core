# Spec Delta

## Purpose

Defines who can read a campaign's details and which fields each role receives, so private campaign data and the invite hash only reach campaign members.

## ADDED Requirements

### Requirement: Campaign details are visible only to campaign members
The system SHALL return a campaign's details only to its DM and to players who have a character in the campaign. Any other authenticated user MUST receive 403.

#### Scenario: DM views campaign details
- **WHEN** the campaign DM requests the campaign details
- **THEN** the system returns 200 with the campaign details

#### Scenario: Player with a character views campaign details
- **WHEN** a user who has a character in the campaign requests the campaign details
- **THEN** the system returns 200 with the campaign details

#### Scenario: Non-member is forbidden
- **WHEN** a user who is neither the DM nor has a character in the campaign requests the campaign details
- **THEN** the system returns 403

#### Scenario: Campaign does not exist
- **WHEN** a user requests the details of a campaign that does not exist
- **THEN** the system returns 404

### Requirement: Invite hash is visible only to the DM
The campaign details response SHALL include the campaign's invite hash only when the requester is the campaign DM. For any other member the invite hash MUST be null.

#### Scenario: DM receives the invite hash
- **WHEN** the DM of a campaign that has an invite requests the campaign details
- **THEN** the response contains the invite hash

#### Scenario: Player does not receive the invite hash
- **WHEN** a player with a character in a campaign that has an invite requests the campaign details
- **THEN** the response invite hash is null
