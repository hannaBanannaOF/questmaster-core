# Spec Delta

## ADDED Requirements

### Requirement: Campaign details include member table names
The campaign details SHALL include `dm_name`, the DM's table name, and each item of `characters` SHALL include `player_name`, its player's table name (see `table-names`). `player_name` is `null` only for a character linked before it was required. Premade characters have no player and carry no name. Only campaign members receive the details; non-members still receive 403.

#### Scenario: Player sees the DM and the party
- **WHEN** a player reads the details of a campaign whose DM is "Ana" and whose characters' players are "Rafa" and "Júlia"
- **THEN** `dm_name` is "Ana" and each character's `player_name` is its player's table name

#### Scenario: Non-member
- **WHEN** a user who is neither the DM nor a player requests the campaign details
- **THEN** the system returns 403 and no names are returned
