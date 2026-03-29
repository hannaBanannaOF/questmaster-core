package character

import (
	characterApp "questmaster-core/internal/character/app"
)

type CreateCharacterUseCase struct {
	r characterApp.CharacterRepository
}

func NewCreateCharacter(r characterApp.CharacterRepository) *CreateCharacterUseCase {
	return &CreateCharacterUseCase{r: r}
}

func (uc *CreateCharacterUseCase) Execute(cmd characterApp.CreateCharacterCommand) (characterApp.CreateCharacterReadModel, error) {
	character, err := uc.r.Create(cmd.Name, cmd.Player, cmd.System, cmd.Hp)
	if err != nil {
		return characterApp.CreateCharacterReadModel{}, err
	}
	return characterApp.MapDomainToCreateCharacterReadModel(character), nil
}
