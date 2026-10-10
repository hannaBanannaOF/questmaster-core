package context

import (
	campaignDomain "questmaster-core/internal/campaign/domain"
	characterDomain "questmaster-core/internal/character/domain"
	"questmaster-core/internal/shared/pagination"
)

func (c *AppContext) SetPage(page pagination.Page) {
	c.Set(string(pageKey), page)
}

func (c *AppContext) Page() pagination.Page {
	v, ok := c.Get(string(pageKey))
	if !ok {
		panic("Page not found in context")
	}

	page, ok := v.(pagination.Page)
	if !ok {
		panic("Page has invalid type")
	}

	return page
}

func (c *AppContext) SetCampaignListFilters(filters campaignDomain.CampaignListFilters) {
	c.Set(string(campaignListFiltersKey), filters)
}

func (c *AppContext) CampaignListFilters() campaignDomain.CampaignListFilters {
	v, ok := c.Get(string(campaignListFiltersKey))
	if !ok {
		panic("CampaignListFilters not found in context")
	}

	filters, ok := v.(campaignDomain.CampaignListFilters)
	if !ok {
		panic("CampaignListFilters has invalid type")
	}

	return filters
}

func (c *AppContext) SetCharacterListFilters(filters characterDomain.CharacterListFilters) {
	c.Set(string(characterListFiltersKey), filters)
}

func (c *AppContext) CharacterListFilters() characterDomain.CharacterListFilters {
	v, ok := c.Get(string(characterListFiltersKey))
	if !ok {
		panic("CharacterListFilters not found in context")
	}

	filters, ok := v.(characterDomain.CharacterListFilters)
	if !ok {
		panic("CharacterListFilters has invalid type")
	}

	return filters
}
