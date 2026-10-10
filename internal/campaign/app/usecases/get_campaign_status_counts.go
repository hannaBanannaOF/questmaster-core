package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
)

type GetCampaignStatusCountsUseCase struct {
	r campaignApp.CampaignRepository
}

func NewGetCampaignStatusCounts(r campaignApp.CampaignRepository) *GetCampaignStatusCountsUseCase {
	return &GetCampaignStatusCountsUseCase{r: r}
}

func (uc *GetCampaignStatusCountsUseCase) Execute(cmd campaignApp.GetCampaignStatusCountsCommand) (campaignApp.CampaignStatusCountsReadModel, error) {
	counts, err := uc.r.CountByStatusForUser(cmd.UserID, cmd.Filters)
	if err != nil {
		return nil, err
	}

	// Every status is present, with 0 when the user has no campaign in it
	result := campaignApp.CampaignStatusCountsReadModel{}
	for _, status := range []campaignDomain.CampaignStatus{
		campaignDomain.StatusDraft,
		campaignDomain.StatusActive,
		campaignDomain.StatusPaused,
		campaignDomain.StatusArchived,
	} {
		result[status.Value()] = counts[status]
	}

	return result, nil
}
