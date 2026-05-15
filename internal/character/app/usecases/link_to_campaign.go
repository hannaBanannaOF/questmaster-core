package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type LinkCharacterToCampaignUseCase struct {
	r characterApp.CharacterRepository
}

func NewLinkCharacterToCampaign(r characterApp.CharacterRepository) *LinkCharacterToCampaignUseCase {
	return &LinkCharacterToCampaignUseCase{
		r: r,
	}
}

func (uc LinkCharacterToCampaignUseCase) LinkToCampaign(campaignID campaignDomain.CampaignID, characterSlug rpgDomain.Slug, userID userDomain.UserID) (characterDomain.Character, error) {
	character, err := uc.r.UpdateCampaign(campaignID, characterSlug, userID)
	if err != nil {
		return characterDomain.Character{}, err
	}

	if character == nil {
		return characterDomain.Character{}, ErrUnavailableCharacter
	}

	return *character, nil
}
