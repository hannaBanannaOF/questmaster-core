package character

import (
	characterApp "questmaster-core/internal/character/app"
	characterDomain "questmaster-core/internal/character/domain"
	"questmaster-core/internal/shared/pagination"
)

type GetCurrentUserCharactersUseCase struct {
	r characterApp.CharacterRepository
}

func NewGetCurrrentUserCharacters(r characterApp.CharacterRepository) *GetCurrentUserCharactersUseCase {
	return &GetCurrentUserCharactersUseCase{r: r}
}

func (uc *GetCurrentUserCharactersUseCase) Execute(cmd characterApp.GetCurrentUserCharactersCommand) (pagination.Result[characterDomain.Character], error) {
	return uc.r.GetAllByPlayerIDWithFilters(cmd.UserID, cmd.Filters, cmd.Page)
}
