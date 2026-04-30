package bootstrap

import (
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignInfra "questmaster-core/internal/campaign/infra/pg"
	"questmaster-core/internal/character"
	"questmaster-core/internal/invite"
	inviteTransport "questmaster-core/internal/invite/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildInviteHandler(db *pgxpool.Pool) *inviteTransport.InviteHandler {
	campaignRepo := campaignInfra.NewCampaignRepositoryPG(db)
	getCampaignFromIDUC := campaignUsecases.NewGetCampaignFromID(campaignRepo)
	characterModule := character.NewCharacterModule(db, getCampaignFromIDUC)
	inviteModule := invite.NewInviteModule(db, getCampaignFromIDUC, characterModule.GetMyCharactersWithoutCampaignUC(), characterModule.LinkCharacterToCampaignUC())

	return inviteTransport.NewInviteHandler(
		inviteModule.GetInviteDetailUC(),
		inviteModule.GetAcceptInviteUC(),
		inviteModule.CreateInviteUC(),
	)
}
