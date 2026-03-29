package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type LinkCharacterToCampaignUseCase struct {
	r characterApp.CharacterRepository
}

func NewLinkCharacterToCampaign(r characterApp.CharacterRepository) *LinkCharacterToCampaignUseCase {
	return &LinkCharacterToCampaignUseCase{
		r: r,
	}
}

func (uc LinkCharacterToCampaignUseCase) LinkToCampaign(campaignID campaignDomain.CampaignID, characterID characterDomain.CharacterID) (characterDomain.Character, error) {
	character, err := uc.r.UpdateCampaign(campaignID, characterID)
	if err != nil {
		return characterDomain.Character{}, err
	}

	if character == nil {
		return characterDomain.Character{}, ErrAlreadyEnrolled
	}

	return *character, nil
}
