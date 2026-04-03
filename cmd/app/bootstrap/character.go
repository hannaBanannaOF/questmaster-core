package bootstrap

import (
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignInfra "questmaster-core/internal/campaign/infra/pg"
	character "questmaster-core/internal/character"
	characterTransport "questmaster-core/internal/character/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildCharacterHandler(db *pgxpool.Pool) *characterTransport.CharactersHandler {
	campaignRepo := campaignInfra.NewCampaignRepositoryPG(db)
	getCampaignFromIDUC := campaignUsecases.NewGetCampaignFromID(campaignRepo)

	module := character.NewCharacterModule(db, getCampaignFromIDUC)

	return characterTransport.NewCharactersHandler(
		module.GetCurrentUserCharactersUC(),
		module.CreateCharacterUC(),
		module.ResolveCharacterSlugUC(),
		module.GetCharacterDetailUC(),
		module.UpdateHPUC(),
		module.DeleteCharacterUC(),
	)
}
