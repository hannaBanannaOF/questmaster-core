package bootstrap

import (
	campaign "questmaster-core/internal/campaign"
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignInfra "questmaster-core/internal/campaign/infra/pg"
	campaignTransport "questmaster-core/internal/campaign/transport/http"

	character "questmaster-core/internal/character"
	invite "questmaster-core/internal/invite"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildCampaignHandler(db *pgxpool.Pool) *campaignTransport.CampaignHandler {
	campaignRepo := campaignInfra.NewCampaignRepositoryPG(db)
	getCampaignFromIDUC := campaignUsecases.NewGetCampaignFromID(campaignRepo)

	characterModule := character.NewCharacterModule(db, getCampaignFromIDUC)
	inviteModule := invite.NewInviteModule(db, getCampaignFromIDUC, characterModule.GetMyCharactersWithoutCampaignUC(), characterModule.LinkCharacterToCampaignUC())

	campaignModule := campaign.NewCampaignModule(
		db,
		characterModule.GetCampaignCharactersUC(),
		inviteModule.GetInviteByCampaignIDUC(),
	)

	return campaignTransport.NewCampaignHandler(
		campaignModule.GetCurrentUserCampaignsUC(),
		campaignModule.ResolveCampaignSlugUC(),
		campaignModule.CreateCampaignUC(),
		campaignModule.UpdateCampaignStatusUC(),
		campaignModule.GetCampaignDetailsUC(),
		campaignModule.DeleteCampaignUC(),
	)
}
