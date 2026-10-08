# Proposal

## Why

Slugs are generated from the name by keeping only ASCII letters and digits. A name with none of them, such as "???" or "東京", produces an empty slug. The application then rejects that slug when reading the row back, so every list that includes the campaign or character fails with a 500 for its owner, permanently.

## What Changes

- When a name produces an empty slug, the slug falls back to `campaign` or `character`, with the usual numeric suffix for uniqueness.
- Existing campaigns and characters with an empty slug get a fallback slug.

## Capabilities

### New Capabilities
- `slugs`: How campaign and character slugs are generated from names, including uniqueness and the fallback for names without ASCII letters or digits.

### Modified Capabilities

None.

## Impact

- **API**: Campaigns and characters with any non-empty name can be created and listed. No contract change.
- **Database**: new migration replacing both slug trigger functions and fixing rows with empty slugs.
