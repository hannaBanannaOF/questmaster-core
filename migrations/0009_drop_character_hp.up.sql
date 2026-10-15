-- Hit points are owned by each game system's service, not by the core (see the remove-character-hp change).
-- This deletes the stored values: back them up before applying if they are still needed.
ALTER TABLE public.character_sheet
    DROP COLUMN IF EXISTS max_hp,
    DROP COLUMN IF EXISTS current_hp;
