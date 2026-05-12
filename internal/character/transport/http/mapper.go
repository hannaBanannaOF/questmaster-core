package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
	userDomain "questmaster-core/internal/user/domain"
)

func convertHP(hp *characterDomain.HP) (*int, *int) {
	var currentHp *int
	var maxHp *int
	if hp != nil {
		savedCurrent := hp.Current()
		savedMax := hp.Max()
		currentHp = &savedCurrent
		maxHp = &savedMax
	}
	return currentHp, maxHp
}

func MapListReadModelToResponse(c characterDomain.Character) CharacterListResponse {
	current, max := convertHP(c.Hp)
	return CharacterListResponse{
		Slug:      c.Slug.Value(),
		Name:      c.Name.Value(),
		System:    c.System.Value(),
		CurrentHP: current,
		MaxHP:     max,
	}
}

func MapDetailReadModelToResponse(c characterDomain.Character) CharacterDetailResponse {
	current, max := convertHP(c.Hp)
	return CharacterDetailResponse{
		Id:        c.Id.Value(),
		Name:      c.Name.Value(),
		System:    c.System.Value(),
		Slug:      c.Slug.Value(),
		MaxHP:     max,
		CurrentHP: current,
	}
}

func MapCurrentHpToResponse(currentHp int) CharacterCurrentHpResponse {
	return CharacterCurrentHpResponse{
		CurrentHP: currentHp,
	}
}

func MapUpdateHPRequestToCommand(req UpdateHPRequest, id characterDomain.CharacterID, userID userDomain.UserID) (characterApp.UpdateHPCommand, error) {
	return characterApp.UpdateHPCommand{
		ID:     id,
		NewHP:  req.NewHP,
		UserID: userID,
	}, nil
}

func MapCreateCharacterReadModelToResponse(rm characterApp.CreateCharacterReadModel) rpgTransport.RpgSlugResponse {
	return rpgTransport.RpgSlugResponse{
		Slug: rm.Slug,
	}
}

func MapResolveSlugReadModelToResponse(rm characterApp.CharacterResolveSlugReadModel) rpgTransport.RpgIdResponse {
	return rpgTransport.RpgIdResponse{
		ID: rm.ID,
	}
}
