# Spec Delta

## Purpose

Defines which campaigns appear in the current user's campaign list and how the player count of a campaign is computed.

## ADDED Requirements

### Requirement: The list contains the user's campaigns as DM and as player
The system SHALL list every campaign where the requester is the DM or has at least one character, each campaign once.

#### Scenario: DM and player campaigns
- **WHEN** a user who is DM of one campaign and has a character in another requests their campaign list
- **THEN** the list contains both campaigns, each once

#### Scenario: Several characters in the same campaign
- **WHEN** a user with two characters in the same campaign requests their campaign list
- **THEN** that campaign appears once

### Requirement: Player count counts distinct players
The player count of a campaign SHALL be the number of distinct players with at least one character in the campaign, regardless of who requests it.

#### Scenario: Player with several characters
- **WHEN** a campaign has two characters of one player and one character of another player
- **THEN** its player count is 2 in the DM's list, in each player's list and in the invite details

#### Scenario: Campaign without characters
- **WHEN** a campaign has no characters
- **THEN** its player count is 0
