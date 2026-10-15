# Spec Delta

## REMOVED Requirements

### Requirement: HP can be updated by the character's player or campaign DM
**Reason**: Hit points are a rule of some game systems; each system-specific service owns its characters' main attribute.
**Migration**: Update HP through the game system's service. `PATCH /character/{characterID}/hp` no longer exists.

### Requirement: Current HP must stay within bounds
**Reason**: The core no longer stores HP.
**Migration**: The game system's service validates its own attribute bounds.

### Requirement: Characters without HP cannot have HP updated
**Reason**: The core no longer stores HP.
**Migration**: None; whether a character has HP is decided by its game system's service.
