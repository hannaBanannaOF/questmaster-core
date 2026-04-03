package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
)

type GetCurrentUserCharactersUseCase struct {
	r characterApp.CharacterRepository
}

func NewGetCurrrentUserCharacters(r characterApp.CharacterRepository) *GetCurrentUserCharactersUseCase {
	return &GetCurrentUserCharactersUseCase{r: r}
}

func (uc *GetCurrentUserCharactersUseCase) Execute(cmd characterApp.GetCurrentUserCharactersCommand) ([]characterDomain.Character, error) {
	characters, err := uc.r.GetAllByPlayerID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}
