DROP TRIGGER IF EXISTS trg_generate_slug ON public.campaign;
DROP FUNCTION IF EXISTS public.generate_unique_campaign_slug();
DROP TABLE IF EXISTS public.campaign;