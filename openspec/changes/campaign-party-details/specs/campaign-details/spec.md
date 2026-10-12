# Spec Delta

## ADDED Requirements

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
