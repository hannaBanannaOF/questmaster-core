# Proposal

## Why

Invites ignore the campaign status: a DM can create an invite for an ARCHIVED campaign, and players can join an ARCHIVED campaign through an existing invite. The DM can also accept their own campaign's invite and join it with one of their characters, which mixes the DM and player roles that the rest of the system keeps separate.

## What Changes

- **BREAKING**: Creating an invite for an ARCHIVED campaign returns 400.
- **BREAKING**: Accepting an invite for an ARCHIVED campaign returns 400.
- **BREAKING**: The campaign DM accepting their own campaign's invite returns 403.
- Documents the existing rules for accepting an invite (own, unlinked character of the campaign's game system).

## Capabilities

### New Capabilities

None.

### Modified Capabilities
- `campaign-invites`: Adds the campaign status rule for creating and accepting invites, the DM restriction on accepting, and the existing character rules for accepting.

## Impact

- **API**: `POST /core/api/v1/invite` and `POST /core/api/v1/invite/{inviteHash}/accept` can return 400 for archived campaigns; accept can return 403 for the DM.
- **Code**: campaign domain policy and errors, create and accept invite use cases, invite module wiring, HTTP error mapping, Swagger annotations.
