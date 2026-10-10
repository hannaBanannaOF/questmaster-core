# Spec Delta

## ADDED Requirements

### Requirement: Invite details tell the requester's role
Any authenticated user holding an invite hash SHALL be able to read the invite details: the campaign slug, name, overview, game system and player count. The response SHALL include `is_dm`, true only when the requester is the DM of the invite's campaign. An unknown invite MUST return 404.

#### Scenario: Player reads the invite
- **WHEN** a user who is not the campaign DM reads the invite details
- **THEN** the system returns 200 with the campaign details
- **AND** `is_dm` is false

#### Scenario: DM reads their own invite
- **WHEN** the campaign DM reads their campaign's invite details
- **THEN** the system returns 200 with the campaign details
- **AND** `is_dm` is true

#### Scenario: Unknown invite
- **WHEN** a user reads the details of an invite hash that does not exist
- **THEN** the system returns 404
