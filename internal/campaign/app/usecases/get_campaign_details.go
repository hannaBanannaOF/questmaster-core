package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
)

type GetCampaignDetailsUseCase struct {
	getCampaignUc           GetCampaignFromIDUseCase
	getCampaignCharactersUc CampaignCharacterFinder
	getCampaignInviteUC     CampaignInviteFinder
}

func NewGetCampaignDetails(
	getCampaignUc GetCampaignFromIDUseCase,
	getCampaignCharactersUc CampaignCharacterFinder,
	getCampaignInviteUC CampaignInviteFinder,
) *GetCampaignDetailsUseCase {
	return &GetCampaignDetailsUseCase{
		getCampaignUc:           getCampaignUc,
		getCampaignCharactersUc: getCampaignCharactersUc,
		getCampaignInviteUC:     getCampaignInviteUC,
	}
}

func (uc *GetCampaignDetailsUseCase) Execute(cmd campaignApp.GetCampaignDetailsCommand) (campaignApp.CampaignDetailsReadModel, error) {
	campaign, err := uc.getCampaignUc.FindByID(cmd.ID)
	if err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	characters, err := uc.getCampaignCharactersUc.GetByCampaignID(cmd.ID)
	if err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	invite, err := uc.getCampaignInviteUC.GetByCampaignID(cmd.ID)
	if err != nil {
		return campaignApp.CampaignDetailsReadModel{}, err
	}

	input := campaignApp.CampaignDetailsInput{
		Campaign:   campaign,
		Characters: characters,
		Invite:     invite,
		UserID:     cmd.UserID,
	}

	return campaignApp.MapDomainToDetailReadModel(input), nil
}
