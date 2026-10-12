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
	if cmd.Name == nil && cmd.MaxHP == nil {
		return characterApp.EditCharacterReadModel{}, ErrNothingToEdit
	}

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

	if err := character.Edit(cmd.UserID, campaignAccess, cmd.Name, cmd.MaxHP); err != nil {
		return characterApp.EditCharacterReadModel{}, err
	}

	saved, err := uc.r.UpdateDetails(*character)
	if err != nil {
		return characterApp.EditCharacterReadModel{}, err
	}

	var currentHP, maxHP *int
	if saved.Hp != nil {
		curr, max := saved.Hp.Current(), saved.Hp.Max()
		currentHP, maxHP = &curr, &max
	}
	return characterApp.EditCharacterReadModel{
		Slug:      saved.Slug.Value(),
		Name:      saved.Name.Value(),
		CurrentHP: currentHP,
		MaxHP:     maxHP,
	}, nil
}
