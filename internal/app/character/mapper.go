package character

import (
	domain "questmaster-core/internal/domain/character"
)

func MapDomainToListReadModel(domain domain.Character) CharacterListReadModel {
	return CharacterListReadModel{
		Slug:   domain.Slug.String(),
		Name:   domain.Name.String(),
		System: string(domain.System),
	}
}

func MapDomainToDetailReadModel(domain domain.Character) CharacterDetailReadModel {
	return CharacterDetailReadModel{
		Name:      domain.Name.String(),
		MaxHp:     domain.Hp.Max(),
		CurrentHp: domain.Hp.Current(),
	}
}
