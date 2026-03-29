package context

import characterDomain "questmaster-core/internal/character/domain"

func (c *AppContext) SetCharacterID(id characterDomain.CharacterID) {
	c.Set(string(characterIDKey), id)
}

func (c *AppContext) CharacterID() characterDomain.CharacterID {
	v, ok := c.Get(string(characterIDKey))
	if !ok {
		panic("CampaignID not found in context")
	}

	id, ok := v.(characterDomain.CharacterID)
	if !ok {
		panic("CampaignID has invalid type")
	}

	return id
}
