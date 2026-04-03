package invite

import (
	characterDomain "questmaster-core/internal/character/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"
)

func MapInviteDetailsReadModelToResponse(rm inviteApp.InviteDetailReadModel) InviteDetailsResponse {
	availableCharacters := make([]InviteDetailsCharacterListItem, 0, len(rm.Characters))

	for _, c := range rm.Characters {
		availableCharacters = append(availableCharacters, InviteDetailsCharacterListItem{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	return InviteDetailsResponse{
		CampaignID:   rm.CampaignID,
		CampaignName: rm.CampaignName,
		Characters:   availableCharacters,
	}
}

func MapAcceptRequestToAcceptCommand(r AcceptInviteRequest, userID userDomain.UserID, hash inviteDomain.InviteHash) inviteApp.AcceptInviteCommand {
	return inviteApp.AcceptInviteCommand{
		Hash:             hash,
		CharacterSheetID: characterDomain.NewCharacterID(r.CharacterSheetID),
		UserID:           userID,
	}
}
