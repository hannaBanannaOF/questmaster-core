package character

import app "questmaster-core/internal/app/character"

func MapListReadModelToResponse(rm app.CharacterListReadModel) CharacterListResponse {
	return CharacterListResponse{
		Slug:   rm.Slug,
		Name:   rm.Name,
		System: rm.System,
	}
}

func MapDetailReadModelToResponse(rm app.CharacterDetailReadModel) CharacterDetailResponse {
	return CharacterDetailResponse{
		Name:      rm.Name,
		MaxHP:     rm.MaxHp,
		CurrentHP: rm.CurrentHp,
	}
}
