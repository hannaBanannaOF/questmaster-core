package campaign

import (
	campaignApp "questmaster-core/internal/campaign/app"
	campaignDomain "questmaster-core/internal/campaign/domain"
	"questmaster-core/internal/shared/pagination"
)

type GetCurrentUserCampaignsUseCase struct {
	r                  campaignApp.CampaignRepository
	myCharactersFinder CampaignPlayerCharactersFinder
}

func NewGetCurrentUserMyCampaigns(r campaignApp.CampaignRepository, myCharactersFinder CampaignPlayerCharactersFinder) *GetCurrentUserCampaignsUseCase {
	return &GetCurrentUserCampaignsUseCase{r: r, myCharactersFinder: myCharactersFinder}
}

func (uc *GetCurrentUserCampaignsUseCase) Execute(cmd campaignApp.GetCurrentUserCampaignsCommand) (pagination.Result[campaignApp.CampaignListItemReadModel], error) {
	page, err := uc.r.ListForUser(cmd.UserID, cmd.Filters, cmd.Page)
	if err != nil {
		return pagination.Result[campaignApp.CampaignListItemReadModel]{}, err
	}

	// The requester's characters are loaded for the returned page only, in one query
	myCharacters := make(map[campaignDomain.CampaignID][]campaignApp.CampaignCharacterRefReadModel, len(page.Items))
	if len(page.Items) > 0 {
		ids := make([]campaignDomain.CampaignID, 0, len(page.Items))
		for _, c := range page.Items {
			ids = append(ids, c.Id)
		}

		characters, err := uc.myCharactersFinder.GetByPlayerInCampaigns(cmd.UserID, ids)
		if err != nil {
			return pagination.Result[campaignApp.CampaignListItemReadModel]{}, err
		}
		for _, ch := range characters {
			if ch.CampaignID == nil {
				continue
			}
			myCharacters[*ch.CampaignID] = append(myCharacters[*ch.CampaignID], campaignApp.CampaignCharacterRefReadModel{
				Slug: ch.Slug.Value(),
				Name: ch.Name.Value(),
			})
		}
	}

	items := make([]campaignApp.CampaignListItemReadModel, 0, len(page.Items))
	for _, c := range page.Items {
		items = append(items, campaignApp.CampaignListItemReadModel{
			Campaign:     c,
			MyCharacters: myCharacters[c.Id],
		})
	}

	return pagination.Result[campaignApp.CampaignListItemReadModel]{Items: items, Total: page.Total}, nil
}
