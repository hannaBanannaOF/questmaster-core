package character

import (
	app "questmaster-core/internal/app/character"
	domain "questmaster-core/internal/domain/character"
	"questmaster-core/internal/domain/rpg"

	"github.com/google/uuid"
)

type CreateCharacterUseCase struct {
	r app.CharacterRepository
}

func NewCreateCharacter(r app.CharacterRepository) *CreateCharacterUseCase {
	return &CreateCharacterUseCase{r: r}
}

func (uc *CreateCharacterUseCase) Execute(name string, system string, hp int, userId uuid.UUID) (string, error) {
	newName, err := domain.NewCharacterName(name)
	if err != nil {
		return "", err
	}

	hpDomain, err := domain.NewHP(hp, hp)
	if err != nil {
		return "", err
	}
	input := app.CreateCharacterInput{
		Name:   newName,
		Hp:     hpDomain,
		Player: rpg.NewUserID(userId),
		System: rpg.System(system),
	}
	newCharacter, err := uc.r.Create(input)
	if err != nil {
		return "", err
	}

	return newCharacter.Slug.String(), nil
}
