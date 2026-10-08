-- Deleting a campaign keeps its characters (unlinked) and removes its invite
ALTER TABLE public.character_sheet
    DROP CONSTRAINT fk_campaign_character_sheet,
    ADD CONSTRAINT fk_campaign_character_sheet FOREIGN KEY (campaign_id)
        REFERENCES public.campaign (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL;

ALTER TABLE public.campaign_invite
    DROP CONSTRAINT fk_campaign_campaign_invite,
    ADD CONSTRAINT fk_campaign_campaign_invite FOREIGN KEY (campaign_id)
        REFERENCES public.campaign (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE;
