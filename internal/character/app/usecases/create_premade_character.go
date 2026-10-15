package character

import (
	characterApp "questmaster-core/internal/character/app"
)

type CreatePremadeCharacterUseCase struct {
	r                   characterApp.CharacterRepository
	getCampaignFromIdUC CharacterCampaingFinder
}

func NewCreatePremadeCharacter(r characterApp.CharacterRepository, getCampaignFromIdUC CharacterCampaingFinder) *CreatePremadeCharacterUseCase {
	return &CreatePremadeCharacterUseCase{r: r, getCampaignFromIdUC: getCampaignFromIdUC}
}

func (uc *CreatePremadeCharacterUseCase) Execute(cmd characterApp.CreatePremadeCharacterCommand) (characterApp.CreateCharacterReadModel, error) {
	campaign, err := uc.getCampaignFromIdUC.FindByID(cmd.CampaignID)
	if err != nil {
		return characterApp.CreateCharacterReadModel{}, err
	}

	if err := campaign.CanCreatePremade(cmd.UserID); err != nil {
		return characterApp.CreateCharacterReadModel{}, err
	}

	character, err := uc.r.CreatePremade(cmd.Name, campaign.Id, campaign.System)
	if err != nil {
		return characterApp.CreateCharacterReadModel{}, err
	}

	return characterApp.MapDomainToCreateCharacterReadModel(character), nil
}
