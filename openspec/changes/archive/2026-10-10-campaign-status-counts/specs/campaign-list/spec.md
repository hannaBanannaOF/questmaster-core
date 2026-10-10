# Spec Delta

## ADDED Requirements

### Requirement: Campaign counts by status
The system SHALL return, for any authenticated user, the number of their own campaigns in each status, with every status present (0 when there is none). The optional `role` filter (`dm` or `player`) restricts the counts like it restricts the campaign list. Campaigns where the requester is neither the DM nor a player MUST NOT be counted. An unknown `role` MUST return 400.

#### Scenario: DM counts
- **WHEN** a user who runs 2 ACTIVE campaigns and 1 ARCHIVED campaign requests the counts with `role=dm`
- **THEN** the system returns 200 with ACTIVE 2, DRAFT 0, PAUSED 0 and ARCHIVED 1

#### Scenario: Both roles
- **WHEN** a user who runs one DRAFT campaign and plays in one ACTIVE campaign requests the counts without `role`
- **THEN** the system returns DRAFT 1 and ACTIVE 1

#### Scenario: Another user's campaigns are not counted
- **WHEN** a user requests the counts
- **THEN** campaigns where the user is neither the DM nor has a character are not counted

#### Scenario: Unknown role
- **WHEN** a user requests the counts with `role=admin`
- **THEN** the system returns 400
