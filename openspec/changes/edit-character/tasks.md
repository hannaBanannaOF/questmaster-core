# Tasks

## 1. Domain

- [ ] 1.1 Add `Character.Edit` (checks `CanUpdate`, renames, `SetMaxHP` clamping the current HP or creating HP); verify with domain tests for player, DM, premade by DM, other user, clamping and a character without HP

## 2. Use case and endpoint

- [ ] 2.1 Add the `EditCharacter` use case (404 unknown, 400 for nothing to change or invalid values, campaign loaded through the existing port); verify with use case tests for each scenario
- [ ] 2.2 Add migration `0008` (character slug trigger also `BEFORE UPDATE OF name`, with its down) and the repository update of `name`, `max_hp` and `current_hp` returning the slug; verify with repository tests that a rename regenerates a unique slug, the old slug no longer resolves, saving the same name keeps the slug, and the down migration restores insert-only
- [ ] 2.3 Add `PATCH /character/{characterID}` with request validation, the 200 response and the Swagger annotation; verify with a handler test

## Workflow follow-up

- Frontend: "Editar" dialog on the character sheet for the player and the campaign DM, redirecting to the new slug after a rename.
- Archive the change after CI passes.
