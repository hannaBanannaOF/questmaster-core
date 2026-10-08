# Design

## Context

See proposal.md for motivation. The delete use case checks `Campaign.CanDelete`, then calls the invite module to delete the invite, then deletes the campaign. Both `character_sheet.campaign_id` and `campaign_invite.campaign_id` are foreign keys with `ON DELETE NO ACTION`.

## Goals / Non-Goals

**Goals:**
- Make the campaign delete a single atomic statement.

**Non-Goals:**
- Changing who can delete or in which statuses (already enforced by `CanDelete`; the spec only documents it).
- A "remove player from campaign" endpoint.

## Decisions

**Foreign key actions instead of application code.** A migration changes the character foreign key to `ON DELETE SET NULL` and the invite foreign key to `ON DELETE CASCADE`, so `DELETE FROM campaign` unlinks characters and removes the invite in the same statement. Alternative: unlink characters and delete the invite from the use case inside a transaction. Rejected because it needs a transaction to span three modules' repositories (the repositories have no transaction support today) and a new port into the character module, for behavior the database expresses in one line each.

**Remove the invite delete use case.** With the cascade, `DeleteInviteUseCase`, `InviteRepository.DeleteByCampaignID` and the `CampaignInviteDeleter` port have no callers. They are removed rather than kept as dead code; a future "revoke invite" feature would need different rules (DM-only, keep campaign) anyway. This also removes the campaign use cases' import of the invite use cases.

## Risks / Trade-offs

- [The rule now lives in the database, invisible from Go code] → The spec documents it, and repository tests against Postgres in CI cover it.
- [Down migration restores `NO ACTION`, which brings the 500 back] → Expected for a rollback.
