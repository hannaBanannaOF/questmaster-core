# Spec Delta

## Purpose

Defines how errors are reported to API clients, so clients get useful status codes without seeing internal details.

## ADDED Requirements

### Requirement: Unexpected errors do not expose internal details
The system SHALL respond to any unexpected error with status 500 and the generic message "Internal server error". The response MUST NOT contain the underlying error text. The underlying error SHALL be logged server-side.

#### Scenario: Database error
- **WHEN** a request fails with an unexpected internal error such as a database error
- **THEN** the system returns 500 with the message "Internal server error"
- **AND** the response does not contain the internal error text
