package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type UpdateCampaignDmNameUseCase struct {
	r campaignApp.CampaignRepository
}

func NewUpdateCampaignDmName(r campaignApp.CampaignRepository) *UpdateCampaignDmNameUseCase {
	return &UpdateCampaignDmNameUseCase{r: r}
}

// Execute changes how the DM is called in the campaign; only the DM may
func (uc *UpdateCampaignDmNameUseCase) Execute(cmd campaignApp.UpdateCampaignDmNameCommand) (rpgDomain.TableName, error) {
	campaign, err := uc.r.FindById(cmd.CampaignID)
	if err != nil {
		return "", err
	}
	if campaign == nil {
		return "", ErrCampaignNotFound
	}

	if err := campaign.CanEdit(cmd.UserID); err != nil {
		return "", err
	}

	if err := uc.r.UpdateDmName(campaign.Id, cmd.DmName); err != nil {
		return "", err
	}
	return cmd.DmName, nil
}
