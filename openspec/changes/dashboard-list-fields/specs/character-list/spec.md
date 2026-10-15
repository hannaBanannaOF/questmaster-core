# Spec Delta

## ADDED Requirements

### Requirement: Each character list item includes its campaign
Each item of the character list SHALL include `campaign`: the `slug` and `name` of the campaign the character is linked to, or `null` when it is linked to none. The list keeps returning only the requester's own characters, so any authenticated user sees the campaign of their own characters only.

#### Scenario: Character in a campaign
- **WHEN** a user whose character "Eleanor" plays in "Innsmouth" requests their character list
- **THEN** the item of "Eleanor" has `campaign` with the slug and name of "Innsmouth"

#### Scenario: Character without a campaign
- **WHEN** a user whose character "Jack" is in no campaign requests their character list
- **THEN** the item of "Jack" has `campaign: null`

#### Scenario: Another player's character is not listed
- **WHEN** a campaign DM requests their own character list
- **THEN** the characters of the campaign's players do not appear, with or without their campaign
