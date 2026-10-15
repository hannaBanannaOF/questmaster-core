# campaign-membership Specification

## Purpose
Defines who can remove a character from a campaign (the player leaving or the campaign DM removing) and what happens to the character and to campaign membership.

## Requirements

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

### Requirement: Only characters with a player in a non-archived campaign can be removed
Removing a character that does not belong to the given campaign, an unclaimed premade character, or any character of an ARCHIVED campaign MUST return 400 and change nothing.

#### Scenario: Character of another campaign
- **WHEN** the DM removes a character that belongs to a different campaign
- **THEN** the system returns 400 and the character stays in its campaign

#### Scenario: Unclaimed premade character
- **WHEN** the DM removes an unclaimed premade character of their campaign through this endpoint
- **THEN** the system returns 400 and the character still exists in the campaign

#### Scenario: Archived campaign
- **WHEN** a player removes their character from an ARCHIVED campaign
- **THEN** the system returns 400 and the character stays in the campaign

### Requirement: A removed character no longer grants membership
After its character is removed, a player with no other character in the campaign SHALL no longer be a campaign member: the campaign leaves their campaign list and its details return 403. The campaign's player count SHALL drop accordingly.

#### Scenario: Last character of a player leaves
- **WHEN** a player with one character in a campaign removes it
- **THEN** the campaign is not in their campaign list and its details return 403 for them
- **AND** the campaign's player count is one lower
