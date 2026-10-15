-- Recreates the columns empty; the values dropped by the up migration are not restored
ALTER TABLE public.character_sheet
    ADD COLUMN IF NOT EXISTS max_hp integer,
    ADD COLUMN IF NOT EXISTS current_hp integer;
