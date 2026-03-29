package context

import inviteDomain "questmaster-core/internal/invite/domain"

func (c *AppContext) SetInviteHash(hash inviteDomain.InviteHash) {
	c.Set(string(inviteHashKey), hash)
}

func (c *AppContext) InviteHash() inviteDomain.InviteHash {
	v, ok := c.Get(string(inviteHashKey))
	if !ok {
		panic("InviteHash not found in context")
	}

	id, ok := v.(inviteDomain.InviteHash)
	if !ok {
		panic("InviteHash has invalid type")
	}

	return id
}
