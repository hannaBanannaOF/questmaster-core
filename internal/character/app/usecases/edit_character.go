package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type EditCharacterUseCase struct {
	r                   characterApp.CharacterRepository
	getCampaignFromIdUC CharacterCampaingFinder
}

func NewEditCharacter(r characterApp.CharacterRepository, getCampaignFromIdUC CharacterCampaingFinder) *EditCharacterUseCase {
	return &EditCharacterUseCase{r: r, getCampaignFromIdUC: getCampaignFromIdUC}
}

func (uc *EditCharacterUseCase) Execute(cmd characterApp.EditCharacterCommand) (characterApp.EditCharacterReadModel, error) {
	character, err := uc.r.FindByID(cmd.ID)
	if err != nil {
		return characterApp.EditCharacterReadModel{}, err
	}
	if character == nil {
		return characterApp.EditCharacterReadModel{}, ErrCharacterNotFound
	}

	var campaignAccess characterDomain.CampaignAccess
	if character.CampaignID != nil {
		campaign, err := uc.getCampaignFromIdUC.FindByID(*character.CampaignID)
		if err != nil {
			return characterApp.EditCharacterReadModel{}, err
		}
		campaignAccess = campaign
	}

	if err := character.Rename(cmd.UserID, campaignAccess, cmd.Name); err != nil {
		return characterApp.EditCharacterReadModel{}, err
	}

	saved, err := uc.r.UpdateDetails(*character)
	if err != nil {
		return characterApp.EditCharacterReadModel{}, err
	}

	return characterApp.EditCharacterReadModel{
		Slug: saved.Slug.Value(),
		Name: saved.Name.Value(),
	}, nil
}
