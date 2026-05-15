package invite

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	userDomain "questmaster-core/internal/user/domain"
)

type GetinviteDetailsCommand struct {
	UserID userDomain.UserID
	Hash   inviteDomain.InviteHash
}

type AcceptInviteCommand struct {
	Hash          inviteDomain.InviteHash
	CharacterSlug rpgDomain.Slug
	UserID        userDomain.UserID
}

type CreateInviteCommand struct {
	CampaignID campaignDomain.CampaignID
}
