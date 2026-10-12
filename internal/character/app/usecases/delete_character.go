package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type DeleteCharacterUseCase struct {
	r                   characterApp.CharacterRepository
	getCampaignFromIdUC CharacterCampaingFinder
}

func NewDeleteCharacter(r characterApp.CharacterRepository, getCampaignFromIdUC CharacterCampaingFinder) *DeleteCharacterUseCase {
	return &DeleteCharacterUseCase{
		r:                   r,
		getCampaignFromIdUC: getCampaignFromIdUC,
	}
}

func (uc *DeleteCharacterUseCase) Execute(cmd characterApp.DeleteCharacterCommand) error {
	character, err := uc.r.FindByID(cmd.ID)
	if err != nil {
		return err
	}

	if character == nil {
		return ErrCharacterNotFound
	}

	// Only an unclaimed premade character needs its campaign, to let the DM delete it
	var campaignAccess characterDomain.CampaignAccess
	if character.IsPremade() && character.CampaignID != nil {
		campaign, err := uc.getCampaignFromIdUC.FindByID(*character.CampaignID)
		if err != nil {
			return err
		}
		campaignAccess = campaign
	}

	if err := character.CanDelete(cmd.UserID, campaignAccess); err != nil {
		return err
	}

	deleted, err := uc.r.DeleteByID(character.Id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrCharacterNotFound
	}

	return nil
}
