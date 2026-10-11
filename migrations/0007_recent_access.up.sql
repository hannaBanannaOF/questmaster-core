-- Last time each user opened a campaign or one of their characters, for the "recent" list sort.
-- One row per user and item; removed with the campaign or character.
CREATE TABLE IF NOT EXISTS public.campaign_access
(
    user_id uuid NOT NULL,
    campaign_id bigint NOT NULL,
    accessed_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT campaign_access_pkey PRIMARY KEY (user_id, campaign_id),
    CONSTRAINT fk_campaign_access_campaign FOREIGN KEY (campaign_id)
        REFERENCES public.campaign (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.character_access
(
    user_id uuid NOT NULL,
    character_id bigint NOT NULL,
    accessed_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT character_access_pkey PRIMARY KEY (user_id, character_id),
    CONSTRAINT fk_character_access_character FOREIGN KEY (character_id)
        REFERENCES public.character_sheet (id)
        ON DELETE CASCADE
);
