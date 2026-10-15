package character

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type UpdatePlayerNameUseCase struct {
	r characterApp.CharacterRepository
}

func NewUpdatePlayerName(r characterApp.CharacterRepository) *UpdatePlayerNameUseCase {
	return &UpdatePlayerNameUseCase{r: r}
}

func (uc *UpdatePlayerNameUseCase) UpdatePlayerName(userID userDomain.UserID, campaignID campaignDomain.CampaignID, name rpgDomain.TableName) (bool, error) {
	return uc.r.UpdatePlayerName(userID, campaignID, name)
}
