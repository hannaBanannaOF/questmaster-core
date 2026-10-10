package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type GetPlayerCharactersInCampaignsUseCase struct {
	r characterApp.CharacterRepository
}

func NewGetPlayerCharactersInCampaigns(r characterApp.CharacterRepository) *GetPlayerCharactersInCampaignsUseCase {
	return &GetPlayerCharactersInCampaignsUseCase{r: r}
}

func (uc *GetPlayerCharactersInCampaignsUseCase) GetByPlayerInCampaigns(
	userID userDomain.UserID,
	campaignIDs []campaignDomain.CampaignID,
) ([]characterDomain.Character, error) {
	return uc.r.GetByPlayerInCampaigns(userID, campaignIDs)
}
