## MODIFIED Requirements

### Requirement: Deleting a campaign keeps its characters and removes its invite
When a campaign is deleted, the system SHALL keep every character linked to it that has a player, unlinked from any campaign and still owned by its player, SHALL delete its unclaimed premade characters, and SHALL remove the campaign's invite. This MUST happen atomically with the campaign deletion.

#### Scenario: Campaign with characters and an invite is deleted
- **WHEN** the campaign DM deletes a campaign that has linked characters and an invite
- **THEN** the system returns 204
- **AND** each character still exists, owned by the same player, with no campaign
- **AND** the invite no longer exists

#### Scenario: Unclaimed premade characters are deleted with the campaign
- **WHEN** the campaign DM deletes a campaign with one claimed and one unclaimed premade character
- **THEN** the claimed character still exists, owned by its player, with no campaign
- **AND** the unclaimed character no longer exists
