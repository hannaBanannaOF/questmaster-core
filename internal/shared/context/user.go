package context

import (
	rpgDomain "questmaster-core/internal/rpg/domain"
)

func (c *AppContext) SetUserID(id rpgDomain.UserID) {
	c.Set(string(userIDKey), id)
}

func (c *AppContext) UserID() rpgDomain.UserID {
	v, ok := c.Get(string(userIDKey))
	if !ok {
		panic("UserID not found in context")
	}

	id, ok := v.(rpgDomain.UserID)
	if !ok {
		panic("UserID has invalid type")
	}

	return id
}
