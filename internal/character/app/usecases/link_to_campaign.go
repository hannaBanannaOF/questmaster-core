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

// LinkToCampaign links the requester own character to the campaign, or makes the requester
// the player of one of the campaign unclaimed premade characters.
func (uc LinkCharacterToCampaignUseCase) LinkToCampaign(campaignID campaignDomain.CampaignID, characterSlug rpgDomain.Slug, userID userDomain.UserID) (characterDomain.Character, error) {
	target, err := uc.r.FindBySlug(characterSlug)
	if err != nil {
		return characterDomain.Character{}, err
	}
	if target != nil && target.IsPremade() {
		return uc.claimPremade(*target, campaignID, userID)
	}

	character, err := uc.r.UpdateCampaign(campaignID, characterSlug, userID)
	if err != nil {
		return characterDomain.Character{}, err
	}

	if character == nil {
		return characterDomain.Character{}, ErrUnavailableCharacter
	}

	return *character, nil
}

func (uc LinkCharacterToCampaignUseCase) claimPremade(premade characterDomain.Character, campaignID campaignDomain.CampaignID, userID userDomain.UserID) (characterDomain.Character, error) {
	if premade.CampaignID == nil || *premade.CampaignID != campaignID {
		return characterDomain.Character{}, ErrUnavailableCharacter
	}

	// Checked first for a clear error; the claim update enforces it again against concurrent claims
	own, err := uc.r.GetByPlayerInCampaigns(userID, []campaignDomain.CampaignID{campaignID})
	if err != nil {
		return characterDomain.Character{}, err
	}
	if len(own) > 0 {
		return characterDomain.Character{}, ErrAlreadyInCampaign
	}

	claimed, err := uc.r.ClaimPremade(premade.Id, campaignID, userID)
	if err != nil {
		return characterDomain.Character{}, err
	}
	if claimed == nil {
		return characterDomain.Character{}, ErrPremadeAlreadyClaimed
	}

	return *claimed, nil
}
