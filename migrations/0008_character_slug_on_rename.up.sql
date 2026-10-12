-- Renaming a character regenerates its slug, with the same function as on insert.
-- Separate from the insert trigger because WHEN can't reference OLD on insert; the WHEN keeps
-- the slug when the same name is saved (recomputing could move it, e.g. "maria-2" to "maria").
CREATE OR REPLACE TRIGGER trg_regenerate_slug_on_rename
BEFORE UPDATE OF name
ON public.character_sheet
FOR EACH ROW
WHEN (OLD.name IS DISTINCT FROM NEW.name)
EXECUTE FUNCTION public.generate_unique_character_sheet_slug();
