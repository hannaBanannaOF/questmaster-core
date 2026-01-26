package character

import (
	app "questmaster-core/internal/app/character"
	domain "questmaster-core/internal/domain/character"
)

type GetCharacterDetailUseCase struct {
	r app.CharacterRepository
}

func NewGetCharacterDetail(r app.CharacterRepository) *GetCharacterDetailUseCase {
	return &GetCharacterDetailUseCase{
		r: r,
	}
}

func (uc *GetCharacterDetailUseCase) Execute(characterId int) (app.CharacterDetailReadModel, error) {
	char, err := uc.r.FindById(domain.CharacterID(characterId))
	if err != nil {
		return app.CharacterDetailReadModel{}, err
	}
	if char == nil {
		return app.CharacterDetailReadModel{}, ErrCharacterNotFound
	}

	return app.MapDomainToDetailReadModel(*char), nil
}
