package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	rpgDomain "questmaster-core/internal/rpg/domain"
	rpgTransport "questmaster-core/internal/rpg/transport/http"
)

func MapListReadModelToResponse(rm characterApp.CharacterListReadModel) CharacterListResponse {
	return CharacterListResponse{
		Slug:   rm.Slug,
		Name:   rm.Name,
		System: rm.System,
	}
}

func MapDetailReadModelToResponse(rm characterApp.CharacterDetailReadModel) CharacterDetailResponse {
	return CharacterDetailResponse{
		Name:      rm.Name,
		MaxHP:     rm.MaxHp,
		CurrentHP: rm.CurrentHp,
	}
}

func MapCurrentHpToResponse(currentHp int) CharacterCurrentHpResponse {
	return CharacterCurrentHpResponse{
		CurrentHP: currentHp,
	}
}

func MapCreateCharacterRequestToCommand(req CreateCharacterRequest, userID rpgDomain.UserID) (characterApp.CreateCharacterCommand, error) {
	name, err := characterDomain.NewCharacterName(req.Name)
	if err != nil {
		return characterApp.CreateCharacterCommand{}, err
	}

	system, err := rpgDomain.NewSystem(req.System)
	if err != nil {
		return characterApp.CreateCharacterCommand{}, nil
	}

	var hp *characterDomain.HP
	if req.Hp != nil {
		h, err := characterDomain.NewHP(*req.Hp, *req.Hp)
		if err != nil {
			return characterApp.CreateCharacterCommand{}, err
		}
		hp = &h
	}

	return characterApp.CreateCharacterCommand{
		Name:   name,
		System: system,
		Hp:     hp,
		Player: userID,
	}, nil
}

func MapUpdateHPRequestToCommand(req UpdateHPRequest, id characterDomain.CharacterID, userID rpgDomain.UserID) (characterApp.UpdateHPCommand, error) {
	hp, err := characterDomain.NewHP(req.NewHP, req.NewHP)
	if err != nil {
		return characterApp.UpdateHPCommand{}, err
	}
	return characterApp.UpdateHPCommand{
		ID:     id,
		NewHP:  hp,
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
