package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type UpdateHPUseCase struct {
	r                   characterApp.CharacterRepository
	getCampaignFromIdUC CharacterCampaingFinder
}

func NewUpdateHP(r characterApp.CharacterRepository, getCampaignFromIdUC CharacterCampaingFinder) *UpdateHPUseCase {
	return &UpdateHPUseCase{
		r:                   r,
		getCampaignFromIdUC: getCampaignFromIdUC,
	}
}

func (uc *UpdateHPUseCase) Execute(cmd characterApp.UpdateHPCommand) (characterApp.UpdateHPReadModel, error) {
	character, err := uc.r.FindByID(cmd.ID)
	if err != nil {
		return characterApp.UpdateHPReadModel{}, err
	}
	if character == nil {
		return characterApp.UpdateHPReadModel{}, ErrCharacterNotFound
	}

	var campaignAccess characterDomain.CampaignAccess

	if character.CampaignID != nil {
		campaign, err := uc.getCampaignFromIdUC.
			FindByID(*character.CampaignID)

		if err != nil {
			return characterApp.UpdateHPReadModel{}, err
		}

		campaignAccess = campaign
	}

	newHP, err := characterDomain.NewHP(cmd.NewHP.Current(), character.Hp.Max())
	if err != nil {
		return characterApp.UpdateHPReadModel{}, err
	}

	if err := character.UpdateHP(newHP, cmd.UserID, campaignAccess); err != nil {
		return characterApp.UpdateHPReadModel{}, err
	}

	newCharacter, err := uc.r.UpdateHP(*character.Hp, character.Id)
	if err != nil {
		return characterApp.UpdateHPReadModel{}, err
	}
	return characterApp.MapDomainToUpdateHPReadModel(newCharacter), nil
}
