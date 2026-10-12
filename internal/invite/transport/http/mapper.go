package invite

import (
	inviteApp "questmaster-core/internal/invite/app"
	inviteDomain "questmaster-core/internal/invite/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

func MapInviteDetailsReadModelToResponse(rm inviteApp.InviteDetailReadModel) InviteDetailsResponse {
	premades := make([]InvitePremadeCharacterResponse, 0, len(rm.PremadeCharacters))
	for _, c := range rm.PremadeCharacters {
		premades = append(premades, InvitePremadeCharacterResponse{Slug: c.Slug, Name: c.Name, CurrentHP: c.CurrentHP, MaxHP: c.MaxHP})
	}

	return InviteDetailsResponse{
		PremadeCharacters:   premades,
		InviteHash:          rm.InviteHash,
		CampaignSlug:        rm.CampaignSlug,
		CampaignName:        rm.CampaignName,
		CampaignOverview:    rm.CampaignOverview,
		CampaignSystem:      rm.CampaignSystem,
		CampaignPlayerCount: rm.CampaignPlayerCount,
		IsDM:                rm.IsDM,
	}
}

func MapAcceptRequestToAcceptCommand(r AcceptInviteRequest, userID userDomain.UserID, hash inviteDomain.InviteHash) (inviteApp.AcceptInviteCommand, error) {
	slug, err := rpgDomain.NewSlug(r.CharacterSlug)
	if err != nil {
		return inviteApp.AcceptInviteCommand{}, err
	}

	return inviteApp.AcceptInviteCommand{
		Hash:          hash,
		CharacterSlug: slug,
		UserID:        userID,
	}, nil
}

func MapInviteDomainToResponse(invite inviteDomain.Invite) InviteCreateResponse {
	return InviteCreateResponse{
		Hash: invite.Hash.Value(),
	}
}
