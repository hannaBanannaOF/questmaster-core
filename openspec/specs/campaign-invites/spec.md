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
