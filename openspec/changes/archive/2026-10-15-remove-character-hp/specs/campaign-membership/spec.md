# Spec Delta

## MODIFIED Requirements

### Requirement: The player or the campaign DM removes a character from a campaign
The system SHALL unlink a character from its campaign when the requester is the character's player or the campaign DM. The character keeps its player and slug, and no longer belongs to any campaign. Any other authenticated user MUST receive 403 and nothing changes. An unknown campaign or character MUST return 404.

#### Scenario: Player leaves the campaign
- **WHEN** a player removes their own character from its ACTIVE campaign
- **THEN** the system returns 204
- **AND** the character still exists, owned by the player, with no campaign

#### Scenario: DM removes a player's character
- **WHEN** the campaign DM removes a player's character from the campaign
- **THEN** the system returns 204
- **AND** the character still exists, owned by its player, with no campaign

#### Scenario: Other member is forbidden
- **WHEN** a player removes another player's character from a campaign they are both in
- **THEN** the system returns 403
- **AND** the character stays in the campaign

#### Scenario: Unknown character
- **WHEN** a user removes a character id that does not exist from a campaign
- **THEN** the system returns 404
