# Spec Delta

## ADDED Requirements

### Requirement: The campaign list shows the requester's characters
Each campaign in the list SHALL include `my_characters`: the requester's characters linked to that campaign, with slug and name, ordered by name. Characters of other players MUST NOT appear in it, even for the campaign DM. It SHALL be empty when the requester has no character in the campaign.

#### Scenario: Player with one character
- **WHEN** a player with one character in a campaign requests their campaign list
- **THEN** that campaign's `my_characters` contains that character's slug and name

#### Scenario: Player with several characters
- **WHEN** a player with two characters in the same campaign requests their campaign list
- **THEN** that campaign appears once, with both characters in `my_characters`

#### Scenario: DM does not see players' characters
- **WHEN** the DM of a campaign with two players requests their campaign list
- **THEN** that campaign's `my_characters` is empty
