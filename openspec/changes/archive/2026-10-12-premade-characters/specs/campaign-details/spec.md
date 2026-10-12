## ADDED Requirements

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
