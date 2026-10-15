## REMOVED Requirements

### Requirement: Members change only their own table name
**Reason**: Table names were edited through `PATCH /campaign/{id}` (DM) and `PATCH /character/{id}` (player), mixing them with campaign and sheet fields.
**Migration**: Use `PATCH /campaign/{id}/table-name` with `{ "name": "..." }`; the requester's role in the campaign decides which name changes.

## ADDED Requirements

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
