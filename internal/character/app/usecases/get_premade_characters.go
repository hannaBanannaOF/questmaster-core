package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type GetPremadeCharactersUseCase struct {
	r characterApp.CharacterRepository
}

func NewGetPremadeCharacters(r characterApp.CharacterRepository) *GetPremadeCharactersUseCase {
	return &GetPremadeCharactersUseCase{r: r}
}

func (uc *GetPremadeCharactersUseCase) GetPremadeByCampaignID(campaignID campaignDomain.CampaignID) ([]characterDomain.Character, error) {
	return uc.r.GetPremadeByCampaignID(campaignID)
}
