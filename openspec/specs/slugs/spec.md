# slugs Specification

## Purpose
Defines how campaign and character slugs are generated from their names, so every campaign and character gets a unique, URL-safe slug.

## Requirements

### Requirement: Slug is derived from the name
The system SHALL generate a slug when a campaign or character is created, by lowercasing the name, removing accents, replacing each run of characters other than ASCII letters and digits with a hyphen, and trimming hyphens from both ends.

#### Scenario: Name with accents and punctuation
- **WHEN** a user creates a campaign named "Mina Assombrada!"
- **THEN** the campaign slug is "mina-assombrada", or "mina-assombrada-N" when that slug is taken

### Requirement: Slugs are unique per entity type
The system SHALL keep slugs unique among campaigns and unique among characters, appending "-2", "-3" and so on when the derived slug is already taken.

#### Scenario: Duplicate name
- **WHEN** a user creates two characters with the same name
- **THEN** the second character's slug is the first one's slug followed by "-2"

### Requirement: Names without ASCII letters or digits get a fallback slug
When the derived slug is empty, the system SHALL use "campaign" for campaigns and "character" for characters as the slug base, with the same uniqueness suffix. Creating and listing such campaigns and characters MUST succeed.

#### Scenario: Campaign named only with symbols
- **WHEN** a user creates a campaign named "???"
- **THEN** the campaign slug is "campaign" or "campaign-N"
- **AND** the campaign appears in the DM's campaign list

#### Scenario: Character named in a non-Latin script
- **WHEN** a user creates a character named "東京"
- **THEN** the character slug is "character" or "character-N"
- **AND** the character appears in the player's character list
