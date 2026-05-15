package context

import (
	characterDomain "questmaster-core/internal/character/domain"
	rpg "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/httperrors"
	"strconv"
)

type QueryFilters struct {
	Filters              map[string]string
	CharacterListFilters *characterDomain.CharacterListFilters
}

func (c *AppContext) SetFilters(filters map[string]string) {
	charFilters := &characterDomain.CharacterListFilters{}

	if v, ok := filters["game_system"]; ok && v != "" {
		newV, err := rpg.NewSystem(v)
		if err != nil {
			c.Error(httperrors.ErrInvalidParam)
			c.Abort()
			return
		}
		charFilters.GameSystem = &newV
	}

	if v, ok := filters["without_campaign"]; ok {
		avail, err := strconv.ParseBool(v)
		if err == nil {
			charFilters.WithoutCampaign = &avail
		}
	}

	c.Set(string(filtersKey), QueryFilters{
		Filters:              filters,
		CharacterListFilters: charFilters,
	})
}

func (c *AppContext) Filters() QueryFilters {
	v, ok := c.Get(string(filtersKey))
	if !ok {
		return QueryFilters{}
	}

	filters, ok := v.(QueryFilters)
	if !ok {
		panic("QueryFilters has invalid type")
	}

	return filters
}
