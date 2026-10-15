# Spec Delta

## ADDED Requirements

### Requirement: The character list can be searched by name

The character list SHALL accept `q` and keep only the requester's characters whose name contains it, ignoring case and accents. Leading and trailing spaces SHALL be trimmed, and an empty `q` SHALL apply no filter. `q` combines with every other filter, sort and page, and `total` reflects it. A `q` longer than 100 characters MUST return 400. Searching never returns characters of other players.

#### Scenario: Search ignores case and accents

- **WHEN** a user with characters "Tomás Ribeiro" and "Jack Brady" requests their character list with `q=TOMAS`
- **THEN** `items` contains only "Tomás Ribeiro" and `total` is 1

#### Scenario: Search combines with filters

- **WHEN** a user requests their character list with `without_campaign=true&q=jack`
- **THEN** `items` contains only their unlinked characters whose name contains "jack"

#### Scenario: Other players' characters are not found

- **WHEN** a campaign DM searches for the name of a player's character
- **THEN** `items` is empty and `total` is 0

#### Scenario: Search too long

- **WHEN** a user requests their character list with a `q` of 101 characters
- **THEN** the system returns 400
