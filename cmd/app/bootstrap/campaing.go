package bootstrap

import (
	usecases "questmaster-core/internal/app/campaign/usecases"
	infra "questmaster-core/internal/infra/campaign/pg"
	transport "questmaster-core/internal/transport/campaign/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildCampaignHandler(db *pgxpool.Pool) *transport.CampaignsHandler {
	repo := infra.NewCampaignRepositoryPG(db)
	fetchUc := usecases.NewFetchMyCampaigns(repo)
	resolveSlugUc := usecases.NewResolveCampaignSlug(repo)
	createCampaignUC := usecases.NewCreateCampaign(repo)
	toggleStatusUC := usecases.NewUpdateStatus(repo)
	return transport.NewCampaignsHandler(fetchUc, resolveSlugUc, createCampaignUC, toggleStatusUC)
}
