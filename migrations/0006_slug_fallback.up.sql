-- Names without ASCII letters or digits (e.g. "???" or "東京") used to produce an empty slug,
-- which the application can't read back. Fall back to the entity type as the slug base.

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

  IF base_slug = '' THEN
    base_slug := 'campaign';
  END IF;

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

  IF base_slug = '' THEN
    base_slug := 'character';
  END IF;

  new_slug := base_slug;

  WHILE EXISTS (SELECT 1 FROM character_sheet WHERE slug = new_slug AND id <> NEW.id) LOOP
    counter := counter + 1;
    new_slug := base_slug || '-' || counter::TEXT;
  END LOOP;

  NEW.slug := new_slug;
  RETURN NEW;
END;
$BODY$;

-- Give existing rows with an empty slug the next free fallback slug
DO $$
DECLARE
  row_id BIGINT;
  new_slug TEXT;
  counter INTEGER;
BEGIN
  FOR row_id IN SELECT id FROM campaign WHERE slug = '' ORDER BY id LOOP
    new_slug := 'campaign';
    counter := 1;
    WHILE EXISTS (SELECT 1 FROM campaign WHERE slug = new_slug) LOOP
      counter := counter + 1;
      new_slug := 'campaign-' || counter::TEXT;
    END LOOP;
    UPDATE campaign SET slug = new_slug WHERE id = row_id;
  END LOOP;

  FOR row_id IN SELECT id FROM character_sheet WHERE slug = '' ORDER BY id LOOP
    new_slug := 'character';
    counter := 1;
    WHILE EXISTS (SELECT 1 FROM character_sheet WHERE slug = new_slug) LOOP
      counter := counter + 1;
      new_slug := 'character-' || counter::TEXT;
    END LOOP;
    UPDATE character_sheet SET slug = new_slug WHERE id = row_id;
  END LOOP;
END;
$$;
