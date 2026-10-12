package character

import (
	"log"

	campaignDomain "questmaster-core/internal/campaign/domain"
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type GetCharacterDetailUseCase struct {
	r                   characterApp.CharacterRepository
	getCampaignFromIdUC CharacterCampaingFinder
}

func NewGetCharacterDetail(r characterApp.CharacterRepository, getCampaignFromIdUC CharacterCampaingFinder) *GetCharacterDetailUseCase {
	return &GetCharacterDetailUseCase{
		r:                   r,
		getCampaignFromIdUC: getCampaignFromIdUC,
	}
}

func (uc *GetCharacterDetailUseCase) Execute(cmd characterApp.GetCharacterDetailsCommand) (characterDomain.Character, error) {
	char, err := uc.r.FindByID(cmd.ID)
	if err != nil {
		return characterDomain.Character{}, err
	}
	if char == nil {
		return characterDomain.Character{}, ErrCharacterNotFound
	}

	// The player needs no lookups; others are checked against the character campaign
	if !char.IsPlayer(cmd.UserID) {
		var campaignAccess characterDomain.CampaignAccess
		isCampaignMember := false
		if char.CampaignID != nil {
			campaign, err := uc.getCampaignFromIdUC.FindByID(*char.CampaignID)
			if err != nil {
				return characterDomain.Character{}, err
			}
			campaignAccess = campaign

			own, err := uc.r.GetByPlayerInCampaigns(cmd.UserID, []campaignDomain.CampaignID{*char.CampaignID})
			if err != nil {
				return characterDomain.Character{}, err
			}
			isCampaignMember = len(own) > 0
		}

		if err := char.CanView(cmd.UserID, campaignAccess, isCampaignMember); err != nil {
			return characterDomain.Character{}, err
		}
	}

	// Only own characters appear in the character list, so only their opens are recorded.
	// Losing one recency update is harmless; failing the read is not
	if char.IsPlayer(cmd.UserID) {
		if err := uc.r.RecordAccess(cmd.UserID, char.Id); err != nil {
			log.Printf("recording access to character %d: %v", char.Id.Value(), err)
		}
	}

	return *char, nil
}
