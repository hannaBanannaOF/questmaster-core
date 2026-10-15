# Spec Delta

## ADDED Requirements

### Requirement: The campaign list can be searched by name

The campaign list SHALL accept `q` and keep only campaigns whose name contains it, ignoring case and accents. Leading and trailing spaces SHALL be trimmed, and an empty `q` SHALL apply no filter. `q` combines with every other filter, sort and page, and `total` reflects it. A `q` longer than 100 characters MUST return 400. Searching never adds campaigns where the requester is neither the DM nor a player.

#### Scenario: Search ignores case and accents

- **WHEN** a user with campaigns "Arkham, 1925" and "Além das Montanhas" requests their campaign list with `q=ALEM`
- **THEN** `items` contains only "Além das Montanhas" and `total` is 1

#### Scenario: Search combines with filters

- **WHEN** a DM requests their campaign list with `role=dm&status=PAUSED&q=boston`
- **THEN** `items` contains only PAUSED campaigns they run whose name contains "boston"

#### Scenario: Other users' campaigns are not found

- **WHEN** a user searches for the name of a campaign where they are neither the DM nor a player
- **THEN** `items` is empty and `total` is 0

#### Scenario: Blank search

- **WHEN** a user requests their campaign list with `q=%20%20`
- **THEN** the result is the same as without `q`

#### Scenario: Search too long

- **WHEN** a user requests their campaign list with a `q` of 101 characters
- **THEN** the system returns 400
