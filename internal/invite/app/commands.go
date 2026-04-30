package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type GetinviteDetailsCommand struct {
	UserID userDomain.UserID
	Hash   inviteDomain.InviteHash
}

type AcceptInviteCommand struct {
	Hash             inviteDomain.InviteHash
	CharacterSheetID characterDomain.CharacterID
	UserID           userDomain.UserID
}

type CreateInviteCommand struct {
	CampaignID campaignDomain.CampaignID
}
