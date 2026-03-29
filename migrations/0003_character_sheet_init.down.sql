DROP TRIGGER IF EXISTS trg_generate_slug ON public.character_sheet;
DROP FUNCTION IF EXISTS public.generate_unique_character_sheet_slug();
DROP TABLE IF EXISTS public.character_sheet;