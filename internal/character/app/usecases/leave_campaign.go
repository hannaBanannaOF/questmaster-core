package character

import (
	characterApp "questmaster-core/internal/character/app"
)

type LeaveCampaignUseCase struct {
	r                   characterApp.CharacterRepository
	getCampaignFromIdUC CharacterCampaingFinder
}

func NewLeaveCampaign(r characterApp.CharacterRepository, getCampaignFromIdUC CharacterCampaingFinder) *LeaveCampaignUseCase {
	return &LeaveCampaignUseCase{r: r, getCampaignFromIdUC: getCampaignFromIdUC}
}

// Execute removes the character from the campaign, by its player (leaving) or the campaign DM (removing)
func (uc *LeaveCampaignUseCase) Execute(cmd characterApp.LeaveCampaignCommand) error {
	character, err := uc.r.FindByID(cmd.CharacterID)
	if err != nil {
		return err
	}
	if character == nil {
		return ErrCharacterNotFound
	}

	campaign, err := uc.getCampaignFromIdUC.FindByID(cmd.CampaignID)
	if err != nil {
		return err
	}

	if character.CampaignID == nil || *character.CampaignID != cmd.CampaignID {
		return ErrNotInCampaign
	}

	if err := character.CanLeaveCampaign(cmd.UserID, campaign); err != nil {
		return err
	}

	// Removed meanwhile by someone else: report it like any character outside the campaign
	unlinked, err := uc.r.UnlinkFromCampaign(character.Id, cmd.CampaignID)
	if err != nil {
		return err
	}
	if !unlinked {
		return ErrNotInCampaign
	}

	return nil
}
