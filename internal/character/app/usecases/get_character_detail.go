package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type GetCharacterDetailUseCase struct {
	r characterApp.CharacterRepository
}

func NewGetCharacterDetail(r characterApp.CharacterRepository) *GetCharacterDetailUseCase {
	return &GetCharacterDetailUseCase{
		r: r,
	}
}

func (uc *GetCharacterDetailUseCase) Execute(characterID characterDomain.CharacterID) (characterApp.CharacterDetailReadModel, error) {
	char, err := uc.r.FindByID(characterID)
	if err != nil {
		return characterApp.CharacterDetailReadModel{}, err
	}
	if char == nil {
		return characterApp.CharacterDetailReadModel{}, ErrCharacterNotFound
	}

	return characterApp.MapDomainToDetailReadModel(*char), nil
}
