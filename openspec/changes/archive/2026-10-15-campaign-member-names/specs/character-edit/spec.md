# Spec Delta

## MODIFIED Requirements

### Requirement: The new name is validated
The name SHALL follow the same rules as character creation, and `player_name` the table name rules (see `table-names`). A request MUST include `name`, `player_name` or both; a request with neither, or with an invalid value, MUST return 400 and change nothing.

#### Scenario: Empty name
- **WHEN** the player sets the character's name to an empty string
- **THEN** the system returns 400 and the name is unchanged

#### Scenario: No name
- **WHEN** the player sends a request with neither `name` nor `player_name`
- **THEN** the system returns 400
