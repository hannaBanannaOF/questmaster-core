package invite

import (
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	inviteInfra "questmaster-core/internal/invite/infra/pg"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InviteModule struct {
	createInviteUC          *inviteUsecases.CreateInviteUseCase
	getInviteByCampaignIDUC *inviteUsecases.GetInviteByCampaignIDUseCase
	getInviteDetailsUC      *inviteUsecases.GetInviteDetailUseCase
	acceptInviteUC          *inviteUsecases.AcceptInviteUseCase
}

func NewInviteModule(
	db *pgxpool.Pool,
	inviteCampaignFinder inviteUsecases.InviteCampaignFinder,
	inviteCharacterCampaignLinker inviteUsecases.InviteCharacterCampaignLinker,
) *InviteModule {
	r := inviteInfra.NewInviteRepositoryPG(db)
	return &InviteModule{
		createInviteUC:          inviteUsecases.NewCreateInvite(r, inviteCampaignFinder),
		getInviteByCampaignIDUC: inviteUsecases.NewGetInviteByCampaignID(r),
		getInviteDetailsUC:      inviteUsecases.NewGetInviteDetail(r, inviteCampaignFinder),
		acceptInviteUC:          inviteUsecases.NewAcceptInvite(r, inviteCampaignFinder, inviteCharacterCampaignLinker),
	}
}

func (m *InviteModule) CreateInviteUC() *inviteUsecases.CreateInviteUseCase {
	return m.createInviteUC
}

func (m *InviteModule) GetInviteByCampaignIDUC() *inviteUsecases.GetInviteByCampaignIDUseCase {
	return m.getInviteByCampaignIDUC
}

func (m *InviteModule) GetInviteDetailUC() *inviteUsecases.GetInviteDetailUseCase {
	return m.getInviteDetailsUC
}

func (m *InviteModule) GetAcceptInviteUC() *inviteUsecases.AcceptInviteUseCase {
	return m.acceptInviteUC
}
