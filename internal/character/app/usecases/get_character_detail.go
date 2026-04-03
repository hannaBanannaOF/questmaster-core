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

func (uc *GetCharacterDetailUseCase) Execute(cmd characterApp.GetCharacterDetailsCommand) (characterDomain.Character, error) {
	char, err := uc.r.FindByID(cmd.ID)
	if err != nil {
		return characterDomain.Character{}, err
	}
	if char == nil {
		return characterDomain.Character{}, ErrCharacterNotFound
	}

	return *char, nil
}
