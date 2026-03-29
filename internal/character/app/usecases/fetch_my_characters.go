package character

import (
	characterApp "questmaster-core/internal/character/app"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type FetchMyCharactersUseCase struct {
	r characterApp.CharacterRepository
}

func NewFetchMyCharacter(r characterApp.CharacterRepository) *FetchMyCharactersUseCase {
	return &FetchMyCharactersUseCase{r: r}
}

func (uc *FetchMyCharactersUseCase) Execute(userID rpgDomain.UserID) ([]characterApp.CharacterListReadModel, error) {
	characters, err := uc.r.GetAllByPlayerID(userID)
	if err != nil {
		return nil, err
	}

	items := make([]characterApp.CharacterListReadModel, 0)

	for _, c := range characters {
		items = append(items, characterApp.MapDomainToListReadModel(c))
	}

	return items, nil
}
