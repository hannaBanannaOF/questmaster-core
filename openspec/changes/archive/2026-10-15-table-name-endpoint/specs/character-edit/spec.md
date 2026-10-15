## MODIFIED Requirements

### Requirement: The new name is validated
The request MUST include `name`, which SHALL follow the same rules as character creation; a missing or invalid name MUST return 400 and change nothing. The character edit does not change table names (see `table-names`); a `player_name` field is ignored.

#### Scenario: Empty name
- **WHEN** the player sets the character's name to an empty string
- **THEN** the system returns 400 and the name is unchanged

#### Scenario: No name
- **WHEN** the player sends a request without `name`
- **THEN** the system returns 400
