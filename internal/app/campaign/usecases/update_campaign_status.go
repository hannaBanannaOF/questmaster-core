package campaign

import (
	app "questmaster-core/internal/app/campaign"
	domain "questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"

	"github.com/google/uuid"
)

type UpdateCampaignStatusUseCase struct {
	r app.CampaignRepository
}

func NewUpdateStatus(r app.CampaignRepository) *UpdateCampaignStatusUseCase {
	return &UpdateCampaignStatusUseCase{r: r}
}

func (uc *UpdateCampaignStatusUseCase) Execute(campaignId int, userId uuid.UUID, newStatus string) (string, error) {
	campaignIdDomain := domain.CampaignID(campaignId)
	newStatusDomain := domain.CampaignStatus(newStatus)
	campaign, err := uc.r.FindById(campaignIdDomain)
	if err != nil {
		return "", err
	}
	if campaign == nil {
		return "", ErrCampaignNotFound
	}

	if err := campaign.ChangeStatus(newStatusDomain, rpg.NewUserID(userId)); err != nil {
		return "", err
	}

	newCampaign, err := uc.r.UpdateStatus(newStatusDomain, campaignIdDomain)
	if err != nil {
		return "", err
	}

	return string(newCampaign.Status), nil
}
