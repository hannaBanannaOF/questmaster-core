package campaign

import (
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignInfra "questmaster-core/internal/campaign/infra/pg"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CampaignModule struct {
	createcampaignUC          *campaignUsecases.CreateCampaignUseCase
	deleteCampaignUC          *campaignUsecases.DeleteCampaignUseCase
	getCurrentUserCampaignsUC *campaignUsecases.GetCurrentUserCampaignsUseCase
	getCampaignDetailUC       *campaignUsecases.GetCampaignDetailsUseCase
	getCampaignFromIDUC       *campaignUsecases.GetCampaignFromIDUseCase
	resolveCampaignSlugUC     *campaignUsecases.ResolveCampaignSlugUseCase
	updateCampaignStatusUC    *campaignUsecases.UpdateCampaignStatusUseCase
}

func NewCampaignModule(
	db *pgxpool.Pool,
	charactersFinder campaignUsecases.CampaignCharacterFinder,
	inviteFinder campaignUsecases.CampaignInviteFinder,
) *CampaignModule {
	r := campaignInfra.NewCampaignRepositoryPG(db)
	getCampaignFromIDUC := campaignUsecases.NewGetCampaignFromID(r)
	return &CampaignModule{
		createcampaignUC:          campaignUsecases.NewCreateCampaign(r),
		deleteCampaignUC:          campaignUsecases.NewDeleteCampaign(r),
		getCurrentUserCampaignsUC: campaignUsecases.NewGetCurrentUserMyCampaigns(r),
		getCampaignDetailUC:       campaignUsecases.NewGetCampaignDetails(*getCampaignFromIDUC, charactersFinder, inviteFinder),
		getCampaignFromIDUC:       getCampaignFromIDUC,
		resolveCampaignSlugUC:     campaignUsecases.NewResolveCampaignSlug(r),
		updateCampaignStatusUC:    campaignUsecases.NewUpdateStatus(r),
	}
}

func (m *CampaignModule) CreateCampaignUC() *campaignUsecases.CreateCampaignUseCase {
	return m.createcampaignUC
}

func (m *CampaignModule) DeleteCampaignUC() *campaignUsecases.DeleteCampaignUseCase {
	return m.deleteCampaignUC
}

func (m *CampaignModule) GetCurrentUserCampaignsUC() *campaignUsecases.GetCurrentUserCampaignsUseCase {
	return m.getCurrentUserCampaignsUC
}

func (m *CampaignModule) GetCampaignDetailsUC() *campaignUsecases.GetCampaignDetailsUseCase {
	return m.getCampaignDetailUC
}

func (m *CampaignModule) GetCampaignFromIDUC() *campaignUsecases.GetCampaignFromIDUseCase {
	return m.getCampaignFromIDUC
}

func (m *CampaignModule) ResolveCampaignSlugUC() *campaignUsecases.ResolveCampaignSlugUseCase {
	return m.resolveCampaignSlugUC
}

func (m *CampaignModule) UpdateCampaignStatusUC() *campaignUsecases.UpdateCampaignStatusUseCase {
	return m.updateCampaignStatusUC
}
