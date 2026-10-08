# Spec Delta

## Purpose

Defines who can delete a campaign, in which statuses, and what happens to the campaign's characters and invite, so players never lose their character sheets.

## ADDED Requirements

### Requirement: Only the campaign DM can delete a campaign
The system SHALL delete a campaign only when the requester is the campaign DM. Any other authenticated user MUST receive 403 and the campaign is kept.

#### Scenario: DM deletes a campaign
- **WHEN** the campaign DM deletes a DRAFT campaign
- **THEN** the system returns 204 and the campaign no longer exists

#### Scenario: Non-DM is forbidden
- **WHEN** a user who is not the campaign DM deletes the campaign
- **THEN** the system returns 403
- **AND** the campaign still exists

#### Scenario: Campaign does not exist
- **WHEN** a user deletes a campaign that does not exist
- **THEN** the system returns 404

### Requirement: Only DRAFT or ARCHIVED campaigns can be deleted
The system SHALL reject deleting an ACTIVE or PAUSED campaign with 403.

#### Scenario: Archived campaign is deleted
- **WHEN** the campaign DM deletes an ARCHIVED campaign
- **THEN** the system returns 204

#### Scenario: Active campaign is rejected
- **WHEN** the campaign DM deletes an ACTIVE campaign
- **THEN** the system returns 403
- **AND** the campaign still exists

### Requirement: Deleting a campaign keeps its characters and removes its invite
When a campaign is deleted, the system SHALL keep every character linked to it, unlinked from any campaign and still owned by its player, and SHALL remove the campaign's invite. This MUST happen atomically with the campaign deletion.

#### Scenario: Campaign with characters and an invite is deleted
- **WHEN** the campaign DM deletes a campaign that has linked characters and an invite
- **THEN** the system returns 204
- **AND** each character still exists, owned by the same player, with no campaign
- **AND** the invite no longer exists
