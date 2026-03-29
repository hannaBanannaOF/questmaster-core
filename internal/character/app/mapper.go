package character

import (
	characterDomain "questmaster-core/internal/character/domain"
)

func MapDomainToListReadModel(domain characterDomain.Character) CharacterListReadModel {
	return CharacterListReadModel{
		Slug:   domain.Slug.Value(),
		Name:   domain.Name.Value(),
		System: domain.System.Value(),
	}
}

func MapDomainToDetailReadModel(domain characterDomain.Character) CharacterDetailReadModel {
	return CharacterDetailReadModel{
		Name:      domain.Name.Value(),
		MaxHp:     domain.Hp.Max(),
		CurrentHp: domain.Hp.Current(),
	}
}

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
