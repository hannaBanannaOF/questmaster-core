package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type GetCampaignCharactersUseCase struct {
	r characterApp.CharacterRepository
}

func NewGetCampaignCharacters(r characterApp.CharacterRepository) *GetCampaignCharactersUseCase {
	return &GetCampaignCharactersUseCase{r: r}
}

func (uc *GetCampaignCharactersUseCase) GetByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return uc.r.GetAllByCampaignID(campaignID)
}
