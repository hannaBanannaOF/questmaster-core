-- Restores the slug functions from migrations 0002 and 0003. Fallback slugs assigned by the
-- up migration are kept, since empty slugs were never valid.

CREATE OR REPLACE FUNCTION public.generate_unique_campaign_slug()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
DECLARE
  base_slug TEXT;
  new_slug TEXT;
  counter INTEGER := 1;
BEGIN
  base_slug := lower(NEW.name);
  base_slug := unaccent(base_slug);
  base_slug := regexp_replace(base_slug, '[^a-z0-9]+', '-', 'g');
  base_slug := trim(both '-' from base_slug);

  new_slug := base_slug;

  WHILE EXISTS (SELECT 1 FROM campaign WHERE slug = new_slug AND id <> NEW.id) LOOP
    counter := counter + 1;
    new_slug := base_slug || '-' || counter::TEXT;
  END LOOP;

  NEW.slug := new_slug;
  RETURN NEW;
END;
$BODY$;

CREATE OR REPLACE FUNCTION public.generate_unique_character_sheet_slug()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
DECLARE
  base_slug TEXT;
  new_slug TEXT;
  counter INTEGER := 1;
BEGIN
  base_slug := lower(NEW.name);
  base_slug := unaccent(base_slug);
  base_slug := regexp_replace(base_slug, '[^a-z0-9]+', '-', 'g');
  base_slug := trim(both '-' from base_slug);

  new_slug := base_slug;

  WHILE EXISTS (SELECT 1 FROM character_sheet WHERE slug = new_slug AND id <> NEW.id) LOOP
    counter := counter + 1;
    new_slug := base_slug || '-' || counter::TEXT;
  END LOOP;

  NEW.slug := new_slug;
  RETURN NEW;
END;
$BODY$;
