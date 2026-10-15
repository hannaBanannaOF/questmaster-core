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
	getCampaignStatusCountsUC *campaignUsecases.GetCampaignStatusCountsUseCase
	updateTableNameUC         *campaignUsecases.UpdateTableNameUseCase
}

func NewCampaignModule(
	db *pgxpool.Pool,
	charactersFinder campaignUsecases.CampaignCharacterFinder,
	inviteFinder campaignUsecases.CampaignInviteFinder,
	myCharactersFinder campaignUsecases.CampaignPlayerCharactersFinder,
	premadeFinder campaignUsecases.CampaignPremadeCharactersFinder,
	playerNames campaignUsecases.CampaignPlayerNameUpdater,
) *CampaignModule {
	r := campaignInfra.NewCampaignRepositoryPG(db)
	getCampaignFromIDUC := campaignUsecases.NewGetCampaignFromID(r)
	return &CampaignModule{
		createcampaignUC:          campaignUsecases.NewCreateCampaign(r),
		deleteCampaignUC:          campaignUsecases.NewDeleteCampaign(r),
		getCurrentUserCampaignsUC: campaignUsecases.NewGetCurrentUserMyCampaigns(r, myCharactersFinder),
		getCampaignDetailUC:       campaignUsecases.NewGetCampaignDetails(*getCampaignFromIDUC, charactersFinder, inviteFinder, premadeFinder, r),
		getCampaignFromIDUC:       getCampaignFromIDUC,
		resolveCampaignSlugUC:     campaignUsecases.NewResolveCampaignSlug(r),
		updateCampaignStatusUC:    campaignUsecases.NewUpdateStatus(r),
		getCampaignStatusCountsUC: campaignUsecases.NewGetCampaignStatusCounts(r),
		updateTableNameUC:         campaignUsecases.NewUpdateTableName(r, playerNames),
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

func (m *CampaignModule) GetCampaignStatusCountsUC() *campaignUsecases.GetCampaignStatusCountsUseCase {
	return m.getCampaignStatusCountsUC
}

func (m *CampaignModule) UpdateTableNameUC() *campaignUsecases.UpdateTableNameUseCase {
	return m.updateTableNameUC
}
