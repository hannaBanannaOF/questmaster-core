# campaign-details Specification

## Purpose
Defines who can read a campaign's details and which fields each role receives, so private campaign data and the invite hash only reach campaign members.

## Requirements

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

### Requirement: Campaign characters carry what a party card needs
Each character in the campaign details SHALL include its `id`, `slug`, `name`, `current_hp`, `max_hp` and `is_mine`, true only when the requester is the character's player. The fields are the same for the DM and for players; non-members still receive 403 for the whole campaign.

#### Scenario: Player sees their own character marked
- **WHEN** a player with one character in the campaign reads the campaign details
- **THEN** that character has `is_mine` true and the other players' characters have `is_mine` false

#### Scenario: DM sees the party
- **WHEN** the campaign DM reads the campaign details of a campaign with two characters
- **THEN** both characters include `slug`, `current_hp` and `max_hp`, with `is_mine` false

#### Scenario: Character without HP
- **WHEN** a character in the campaign has no HP set
- **THEN** its `current_hp` and `max_hp` are null

### Requirement: The DM sees the campaign's unclaimed premade characters
The campaign details `characters` SHALL contain only characters with a player. The response SHALL include `premade_characters`, the campaign's unclaimed premade characters ordered by name, each with `id`, `slug`, `name`, `current_hp` and `max_hp`, only when the requester is the campaign DM. For any other member it MUST be an empty list. A premade character does not count as a player nor make anyone a campaign member.

#### Scenario: DM sees premade characters apart from the party
- **WHEN** the campaign DM reads the details of a campaign with one player character and one unclaimed premade character
- **THEN** `characters` contains only the player character
- **AND** `premade_characters` contains the premade character

#### Scenario: Player does not see premade characters
- **WHEN** a player with a character in the campaign reads the details of a campaign with an unclaimed premade character
- **THEN** `premade_characters` is empty

#### Scenario: Premade characters do not count as players
- **WHEN** a campaign has one player character and two unclaimed premade characters
- **THEN** its player count is 1
