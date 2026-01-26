package campaign

import (
	app "questmaster-core/internal/app/campaign"
	domain "questmaster-core/internal/domain/campaign"
	"questmaster-core/internal/domain/rpg"

	"github.com/google/uuid"
)

type CreateCampaignUseCase struct {
	r app.CampaignRepository
}

func NewCreateCampaign(r app.CampaignRepository) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{r: r}
}

func (uc *CreateCampaignUseCase) Execute(name string, overview *string, system string, dmId uuid.UUID) (string, error) {
	newName, err := domain.NewCampaignName(name)
	if err != nil {
		return "", err
	}
	var campaignOverview *domain.CampaignOverview

	if overview != nil {
		o, err := domain.NewCampaignOverview(*overview)
		if err != nil {
			return "", nil
		}
		campaignOverview = &o
	}
	input := app.CreateCampaignInput{
		Name:     newName,
		Overview: campaignOverview,
		Dm:       rpg.NewUserID(dmId),
		System:   rpg.System(system),
	}
	newCamapign, err := uc.r.Create(input)
	if err != nil {
		return "", err
	}

	return newCamapign.Slug.String(), nil
}
