package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaigndomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/pagination"
)

type GetCurrentUserCampaignsUseCase struct {
	r campaignApp.CampaignRepository
}

func NewGetCurrentUserMyCampaigns(r campaignApp.CampaignRepository) *GetCurrentUserCampaignsUseCase {
	return &GetCurrentUserCampaignsUseCase{r: r}
}

func (uc *GetCurrentUserCampaignsUseCase) Execute(cmd campaignApp.GetCurrentUserCampaignsCommand) (pagination.Result[campaigndomain.Campaign], error) {
	return uc.r.ListForUser(cmd.UserID, cmd.Filters, cmd.Page)
}
