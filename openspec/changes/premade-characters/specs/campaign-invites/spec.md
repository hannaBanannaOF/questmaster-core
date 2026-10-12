## MODIFIED Requirements

### Requirement: Accepting an invite links an eligible character
Accepting an invite SHALL link the given character to the invite's campaign only if the character belongs to the requester, is not linked to any campaign, and has the campaign's game system. Accepting with an unclaimed premade character of the invite's campaign SHALL instead make the requester its player. Otherwise the system MUST return 400 and change nothing. A premade character claimed by someone else in the meantime MUST return 409 and change nothing. An unknown invite MUST return 404.

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

#### Scenario: Claiming a premade character
- **WHEN** a player accepts an invite with an unclaimed premade character of the invite's campaign
- **THEN** the system returns 204
- **AND** the requester is the character's player and the character stays in the campaign

#### Scenario: Premade character already claimed
- **WHEN** two players accept the invite with the same premade character
- **THEN** the first receives 204 and becomes its player
- **AND** the second receives 409 and the character keeps the first player

#### Scenario: Premade character of another campaign
- **WHEN** a player accepts an invite with an unclaimed premade character of a different campaign
- **THEN** the system returns 400 and nothing changes

### Requirement: Invite details tell the requester's role
Any authenticated user holding an invite hash SHALL be able to read the invite details: the campaign slug, name, overview, game system and player count. The response SHALL include `is_dm`, true only when the requester is the DM of the invite's campaign, and `premade_characters`, the campaign's unclaimed premade characters ordered by name, each with `slug`, `name`, `current_hp` and `max_hp`. An unknown invite MUST return 404.

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

#### Scenario: Invite lists unclaimed premade characters
- **WHEN** a user reads the invite of a campaign with one unclaimed and one claimed premade character
- **THEN** `premade_characters` contains only the unclaimed one, with its slug, name and HP
