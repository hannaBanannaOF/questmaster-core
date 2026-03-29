package context

import rpgDomain "questmaster-core/internal/rpg/domain"

func (c *AppContext) SetSlug(slug rpgDomain.Slug) {
	c.Set(string(slugKey), slug)
}

func (c *AppContext) Slug() rpgDomain.Slug {
	v, ok := c.Get(string(slugKey))
	if !ok {
		panic("Slug not found in context")
	}

	slug, ok := v.(rpgDomain.Slug)
	if !ok {
		panic("Slug has invalid type")
	}

	return slug
}
