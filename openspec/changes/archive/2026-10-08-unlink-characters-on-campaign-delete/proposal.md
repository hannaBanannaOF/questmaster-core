# Proposal

## Why

Deleting a DRAFT or ARCHIVED campaign that has characters fails with a 500, because the character foreign key blocks the delete. Deleting the invite and then the campaign also runs as two separate statements, so a failure in between leaves the campaign without its invite.

## What Changes

- Deleting a campaign unlinks its characters (they stay with their players, without a campaign) instead of failing.
- The campaign's invite is removed in the same atomic delete.
- The separate "delete invite by campaign" use case and repository method are removed, since nothing else uses them.

## Capabilities

### New Capabilities
- `campaign-deletion`: Who can delete a campaign, in which statuses, and what happens to its characters and invite.

### Modified Capabilities

None.

## Impact

- **API**: `DELETE /core/api/v1/campaign/{campaignID}` returns 204 for campaigns with characters instead of 500. No contract change.
- **Database**: new migration changing the `campaign_id` foreign keys of `character_sheet` (ON DELETE SET NULL) and `campaign_invite` (ON DELETE CASCADE).
- **Code**: campaign delete use case and module wiring, invite module (delete use case removed), bootstrap.
