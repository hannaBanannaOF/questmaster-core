# Spec Delta

## ADDED Requirements

### Requirement: Campaign characters identify a party card
Each character in the campaign details SHALL include its `id`, `slug`, `name` and `is_mine`, true only when the requester is the character's player. The fields are the same for the DM and for players; non-members still receive 403 for the whole campaign.

#### Scenario: Player sees their own character marked
- **WHEN** a player with one character in the campaign reads the campaign details
- **THEN** that character has `is_mine` true and the other players' characters have `is_mine` false

#### Scenario: DM sees the party
- **WHEN** the campaign DM reads the campaign details of a campaign with two characters
- **THEN** both characters include `slug` and `name`, with `is_mine` false

## MODIFIED Requirements

### Requirement: The DM sees the campaign's unclaimed premade characters
The campaign details `characters` SHALL contain only characters with a player. The response SHALL include `premade_characters`, the campaign's unclaimed premade characters ordered by name, each with `id`, `slug` and `name`, only when the requester is the campaign DM. For any other member it MUST be an empty list. A premade character does not count as a player nor make anyone a campaign member.

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

## REMOVED Requirements

### Requirement: Campaign characters carry what a party card needs
**Reason**: Its HP fields and the "Character without HP" scenario no longer apply.
**Migration**: Replaced by "Campaign characters identify a party card", without `current_hp` and `max_hp`.
