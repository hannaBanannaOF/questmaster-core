package character

import (
	app "questmaster-core/internal/app/character"
	"questmaster-core/internal/domain/rpg"

	"github.com/google/uuid"
)

type FetchMyCharactersUseCase struct {
	r app.CharacterRepository
}

func NewFetchMyCharacter(r app.CharacterRepository) *FetchMyCharactersUseCase {
	return &FetchMyCharactersUseCase{r: r}
}

func (uc *FetchMyCharactersUseCase) Execute(userId uuid.UUID) ([]app.CharacterListReadModel, error) {
	characters, err := uc.r.GetAllByPlayerId(rpg.NewUserID(userId))
	if err != nil {
		return nil, err
	}

	items := make([]app.CharacterListReadModel, 0)

	for _, c := range characters {
		items = append(items, app.MapDomainToListReadModel(c))
	}

	return items, nil
}
