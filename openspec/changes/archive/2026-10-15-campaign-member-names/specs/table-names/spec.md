# Spec Delta

## Purpose

Defines how campaign members name themselves at a table: the name each one gives for a single campaign, when it is set, who can change it and when it is removed.

## ADDED Requirements

### Requirement: The DM names themselves when creating a campaign
Creating a campaign SHALL require `dm_name`, the name the DM is shown with in that campaign. It is trimmed and MUST have 1 to 50 characters; otherwise the system returns 400 and no campaign is created.

#### Scenario: Campaign created with the DM's table name
- **WHEN** a user creates a campaign with `dm_name` "Mestre Ana"
- **THEN** the system returns 201 and the campaign's DM is shown as "Mestre Ana"

#### Scenario: Missing DM name
- **WHEN** a user creates a campaign without `dm_name`, or with only spaces
- **THEN** the system returns 400 and no campaign is created

### Requirement: A player names themselves when accepting an invite
Accepting an invite SHALL require `player_name`, the name the player is shown with in that campaign, for their own character and for a claimed premade character alike. It is trimmed and MUST have 1 to 50 characters; otherwise the system returns 400 and nothing changes.

#### Scenario: Player joins with a table name
- **WHEN** a player accepts an invite with their character and `player_name` "Marina"
- **THEN** the character joins the campaign and its player is shown there as "Marina"

#### Scenario: Missing player name
- **WHEN** a player accepts an invite without `player_name`
- **THEN** the system returns 400 and the character is not linked

### Requirement: Table names are removed with the table
A player's table name SHALL be cleared when their character leaves or is removed from the campaign, and when the campaign is deleted. Deleting the campaign or the character deletes the names stored on them.

#### Scenario: Player leaves and rejoins
- **WHEN** a player who joined as "Marina" leaves the campaign and later accepts its invite again with `player_name` "Mari"
- **THEN** they are shown as "Mari"

#### Scenario: Campaign deleted
- **WHEN** the DM deletes a campaign
- **THEN** the characters that were in it no longer carry a table name

### Requirement: Members change only their own table name
The campaign DM SHALL be able to change `dm_name`, and a character's player SHALL be able to change their `player_name` while the character is in a campaign, with the same length rule. Nobody can change another member's table name. Changing `player_name` on a character with no campaign MUST return 400.

#### Scenario: DM renames themselves
- **WHEN** the campaign DM changes `dm_name` to "Ana"
- **THEN** the system returns 200 and the DM is shown as "Ana"

#### Scenario: Player renames themselves
- **WHEN** a player changes their character's `player_name` to "Mari"
- **THEN** the system returns 200 and they are shown as "Mari"

#### Scenario: DM changes a player's name
- **WHEN** the campaign DM sends `player_name` for a player's character
- **THEN** the system returns 403 and the name is unchanged

#### Scenario: Player changes the DM name
- **WHEN** a player sends `dm_name` for the campaign
- **THEN** the system returns 403

#### Scenario: Character without a campaign
- **WHEN** a player sends `player_name` for a character that is in no campaign
- **THEN** the system returns 400
