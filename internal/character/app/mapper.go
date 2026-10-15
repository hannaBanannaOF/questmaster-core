package character

import (
	characterDomain "questmaster-core/internal/character/domain"
)

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
