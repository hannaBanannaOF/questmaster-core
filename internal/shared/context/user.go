package context

import (
	userDomain "questmaster-core/internal/user/domain"
)

func (c *AppContext) SetUser(data userDomain.User) {
	c.Set(string(userKey), data)
}

func (c *AppContext) User() userDomain.User {
	v, ok := c.Get(string(userKey))
	if !ok {
		panic("User not found in context")
	}

	data, ok := v.(userDomain.User)
	if !ok {
		panic("User has invalid type")
	}

	return data
}

func (c *AppContext) UserID() userDomain.UserID {
	v, ok := c.Get(string(userKey))
	if !ok {
		panic("User not found in context")
	}

	data, ok := v.(userDomain.User)
	if !ok {
		panic("User has invalid type")
	}

	return data.Id
}
