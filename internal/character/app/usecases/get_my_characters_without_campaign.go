package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type GetMyCharactersWithoutCampaignUseCase struct {
	r characterApp.CharacterRepository
}

func NewMyGetCharactersWithoutCampaign(r characterApp.CharacterRepository) *GetMyCharactersWithoutCampaignUseCase {
	return &GetMyCharactersWithoutCampaignUseCase{
		r: r,
	}
}

func (uc *GetMyCharactersWithoutCampaignUseCase) GetBySystemAndCampaignIDNull(userID userDomain.UserID, system rpgDomain.System) ([]characterDomain.Character, error) {
	return uc.r.GetAllByUserIDAndCampaignIDNullAndSystem(userID, system)
}
