# user-profile Specification

## Purpose
Defines the authenticated user's profile, which is built from the identity token claims rather than stored by this service.

## Requirements

### Requirement: Profile is built from identity token claims
The system SHALL return the authenticated user's id, username, name and surname from the identity token. Name and surname MUST be null when the token has no given name.

#### Scenario: Token has a given name
- **WHEN** a user whose token has given and family names requests their profile
- **THEN** the system returns 200 with the name and surname from the token

#### Scenario: Token has no given name
- **WHEN** a user whose token has no given name requests their profile
- **THEN** the system returns 200 with null name and null surname
