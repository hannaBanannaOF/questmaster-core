package invite

import (
	characterDomain "questmaster-core/internal/character/domain"
	inviteDomain "questmaster-core/internal/invite/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type GetinviteDetailsCommand struct {
	UserID rpgDomain.UserID
	Hash   inviteDomain.InviteHash
}

type AcceptInviteCommand struct {
	Hash             inviteDomain.InviteHash
	CharacterSheetID characterDomain.CharacterID
	UserID           rpgDomain.UserID
}
