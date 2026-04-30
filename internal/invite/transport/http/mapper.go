package invite

import (
	characterDomain "questmaster-core/internal/character/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"
)

func MapInviteDetailsReadModelToResponse(rm inviteApp.InviteDetailReadModel) InviteDetailsResponse {
	return InviteDetailsResponse{
		CampaignID:          rm.CampaignID,
		CampaignName:        rm.CampaignName,
		CampaignOverview:    rm.CampaignOverview,
		CampaignSystem:      rm.CampaignSystem,
		CampaignPlayerCount: rm.CampaignPlayerCount,
	}
}

func MapAcceptRequestToAcceptCommand(r AcceptInviteRequest, userID userDomain.UserID, hash inviteDomain.InviteHash) inviteApp.AcceptInviteCommand {
	return inviteApp.AcceptInviteCommand{
		Hash:             hash,
		CharacterSheetID: characterDomain.NewCharacterID(r.CharacterSheetID),
		UserID:           userID,
	}
}

func MapInviteDomainToResponse(invite inviteDomain.Invite) InviteCreateResponse {
	return InviteCreateResponse{
		Hash: invite.Hash.Value(),
	}
}
