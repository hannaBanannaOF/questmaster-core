## MODIFIED Requirements

### Requirement: Slug is derived from the name
The system SHALL generate a slug when a campaign or character is created, and regenerate a character's slug when its name changes, by lowercasing the name, removing accents, replacing each run of characters other than ASCII letters and digits with a hyphen, and trimming hyphens from both ends.

#### Scenario: Name with accents and punctuation
- **WHEN** a user creates a campaign named "Mina Assombrada!"
- **THEN** the campaign slug is "mina-assombrada", or "mina-assombrada-N" when that slug is taken

#### Scenario: Character renamed
- **WHEN** a player renames a character to "Jão Ninguém"
- **THEN** the character slug becomes "jao-ninguem", or "jao-ninguem-N" when that slug is taken
