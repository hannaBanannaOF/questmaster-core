package character

import (
	characterDomain "questmaster-core/internal/character/domain"
)

func MapDomainToUpdateHPReadModel(character characterDomain.Character) UpdateHPReadModel {
	return UpdateHPReadModel{
		CurrentHP: character.Hp.Current(),
	}
}

func MapDomainToCharacterResolveSlugReadModel(character characterDomain.Character) CharacterResolveSlugReadModel {
	return CharacterResolveSlugReadModel{
		ID: character.Id.Value(),
	}
}

func MapDomainToCreateCharacterReadModel(character characterDomain.Character) CreateCharacterReadModel {
	return CreateCharacterReadModel{
		Slug: character.Slug.Value(),
	}
}
