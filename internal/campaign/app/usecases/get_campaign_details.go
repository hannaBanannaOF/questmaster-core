package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
)

type GetCampaignDetailsUseCase struct {
	getCampaignUc   GetCampaignFromIDUseCase
	getCharactersUc CampaignCharacterFinder
}

func NewGetCampaignDetails(
	getCampaignUc GetCampaignFromIDUseCase,
	getCharactersUc CampaignCharacterFinder,
) *GetCampaignDetailsUseCase {
	return &GetCampaignDetailsUseCase{
		getCampaignUc:   getCampaignUc,
		getCharactersUc: getCharactersUc,
	}
}

func (uc *GetCampaignDetailsUseCase) Execute(cmd campaignApp.GetCampaignDetailsCommand) (campaignApp.CampaignDetailsReadModel, error) {
	campaign, err := uc.getCampaignUc.FindByID(cmd.ID)
	if err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	characters, err := uc.getCharactersUc.GetByCampaignID(cmd.ID)
	if err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	return campaignApp.MapDomainToDetailReadModel(campaign, characters, cmd.UserID), nil
}
