# Design

## Context

See proposal.md (Why). HP lives in the character domain (`HP` value object with `current` and `max`, `Character.Hp`), in `character_sheet.max_hp/current_hp`, in the update HP use case and endpoint, in character, premade and edit use cases, and in every response that shows characters. `Character.CanUpdate` (player or campaign DM) is shared by the HP update and the character edit.

## Goals / Non-Goals

**Goals:**
- No HP anywhere in the core: domain, API, database.

**Non-Goals:**
- Moving the existing HP values to the system-specific services. They own their attributes from now on; migrating data, if wanted, happens before deploying this change, outside the core.
- Replacing HP with a generic "main attribute" in the core.

## Decisions

**Drop the columns instead of leaving them unused.** Unused columns invite new code to read them again, and `cs.*` queries would keep loading them. The down migration recreates them nullable and empty, so rolling back the code works, but the values are gone. Alternative: keep the columns for a release as a backup. Rejected: the services start from their own data; a database backup taken before the deploy covers recovery.

**Character edit keeps its endpoint and permission.** `PATCH /character/{id}` stays, with `name` only, still allowed for the player and the campaign DM through `CanUpdate`. A request without a valid name returns 400, through the existing name validation.

**Remove, don't deprecate, the API fields.** The frontend is the only client and ships together with this change; keeping `current_hp: null` would only postpone the cleanup.

## Risks / Trade-offs

- [HP values are lost on deploy] → Called out as BREAKING (data); take a database backup, or export the columns, before running migration 0009.
- [Clients still sending `hp` or `max_hp`] → Unknown JSON fields are ignored; requests keep working without HP.
