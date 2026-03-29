package campaign

import (
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignInfra "questmaster-core/internal/campaign/infra/pg"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CampaignModule struct {
	createcampaignUC            *campaignUsecases.CreateCampaignUseCase
	deleteCampaignUC            *campaignUsecases.DeleteCampaignUseCase
	fetchMyCampaingUC           *campaignUsecases.FetchMyCampaignsUseCase
	getCampaignDetailUC         *campaignUsecases.GetCampaignDetailsUseCase
	getCampaignFromIDUC         *campaignUsecases.GetCampaignFromIDUseCase
	getOrCreateCampaignInviteUC *campaignUsecases.GetOrCreateCampaignInviteUseCase
	resolveCampaignSlugUC       *campaignUsecases.ResolveCampaignSlugUseCase
	updateCampaignStatusUC      *campaignUsecases.UpdateCampaignStatusUseCase
}

func NewCampaignModule(
	db *pgxpool.Pool,
	charactersFinder campaignUsecases.CampaignCharacterFinder,
	campaignInviteFinder campaignUsecases.CampaignInviteFinder,
	campaignInviteCreator campaignUsecases.CampaignInviteCreator,
) *CampaignModule {
	r := campaignInfra.NewCampaignRepositoryPG(db)
	getCampaignFromIDUC := campaignUsecases.NewGetCampaignFromID(r)
	return &CampaignModule{
		createcampaignUC:            campaignUsecases.NewCreateCampaign(r),
		deleteCampaignUC:            campaignUsecases.NewDeleteCampaign(r),
		fetchMyCampaingUC:           campaignUsecases.NewFetchMyCampaigns(r),
		getCampaignDetailUC:         campaignUsecases.NewGetCampaignDetails(*getCampaignFromIDUC, charactersFinder),
		getCampaignFromIDUC:         getCampaignFromIDUC,
		getOrCreateCampaignInviteUC: campaignUsecases.NewGetOrCreateCampaignInvite(*getCampaignFromIDUC, campaignInviteFinder, campaignInviteCreator),
		resolveCampaignSlugUC:       campaignUsecases.NewResolveCampaignSlug(r),
		updateCampaignStatusUC:      campaignUsecases.NewUpdateStatus(r),
	}
}

func (m *CampaignModule) CreateCampaignUC() *campaignUsecases.CreateCampaignUseCase {
	return m.createcampaignUC
}

func (m *CampaignModule) DeleteCampaignUC() *campaignUsecases.DeleteCampaignUseCase {
	return m.deleteCampaignUC
}

func (m *CampaignModule) FetchMyCampaignsUC() *campaignUsecases.FetchMyCampaignsUseCase {
	return m.fetchMyCampaingUC
}

func (m *CampaignModule) GetCampaignDetailsUC() *campaignUsecases.GetCampaignDetailsUseCase {
	return m.getCampaignDetailUC
}

func (m *CampaignModule) GetCampaignFromIDUC() *campaignUsecases.GetCampaignFromIDUseCase {
	return m.getCampaignFromIDUC
}

func (m *CampaignModule) GetOrCreateCampaignInviteUC() *campaignUsecases.GetOrCreateCampaignInviteUseCase {
	return m.getOrCreateCampaignInviteUC
}

func (m *CampaignModule) ResolveCampaignSlugUC() *campaignUsecases.ResolveCampaignSlugUseCase {
	return m.resolveCampaignSlugUC
}

func (m *CampaignModule) UpdateCampaignStatusUC() *campaignUsecases.UpdateCampaignStatusUseCase {
	return m.updateCampaignStatusUC
}
