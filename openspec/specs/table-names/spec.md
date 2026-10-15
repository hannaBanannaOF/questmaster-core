# table-names Specification

## Purpose
Defines how campaign members name themselves at a table: the name each one gives for a single campaign, when it is set, who can change it and when it is removed.

## Requirements

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

### Requirement: Members change their own table name in the campaign
`PATCH /campaign/{id}/table-name` SHALL change the requester's own table name in that campaign: the DM's `dm_name`, or the `player_name` of the requester's character in the campaign. `name` is trimmed and MUST have 1 to 50 characters; otherwise the system returns 400 and nothing changes. A user who is neither the DM nor a player in the campaign MUST get 403, and an unknown campaign 404. It works in every campaign status.

#### Scenario: DM renames themselves
- **WHEN** the campaign DM sends `name` "Ana"
- **THEN** the system returns 200 with the name and the DM is shown as "Ana"

#### Scenario: Player renames themselves
- **WHEN** a player with a character in the campaign sends `name` "Mari"
- **THEN** the system returns 200 with the name and they are shown as "Mari", while the other members' names are unchanged

#### Scenario: Outsider
- **WHEN** a user with no character in the campaign, who is not its DM, sends a name
- **THEN** the system returns 403 and no name changes

#### Scenario: Invalid name
- **WHEN** a member sends an empty name, only spaces, or more than 50 characters
- **THEN** the system returns 400 and the name is unchanged

#### Scenario: Unknown campaign
- **WHEN** a user sends a name for a campaign that does not exist
- **THEN** the system returns 404

### Requirement: The campaign has no generic edit route
`PATCH /campaign/{id}` SHALL NOT exist; status changes keep their own route and table names use the table-name route.

#### Scenario: Generic campaign patch
- **WHEN** a client sends `PATCH /campaign/{id}`
- **THEN** the system returns 404
