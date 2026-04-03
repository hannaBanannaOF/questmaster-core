package character

import (
	characterApp "questmaster-core/internal/character/app"
	rpgApp "questmaster-core/internal/rpg/app"
)

type ResolveCharacterSlugUseCase struct {
	r characterApp.CharacterRepository
}

func NewResolveCharacterSlug(r characterApp.CharacterRepository) *ResolveCharacterSlugUseCase {
	return &ResolveCharacterSlugUseCase{r: r}
}

func (uc *ResolveCharacterSlugUseCase) Execute(cmd rpgApp.ResolveSlugCommand) (characterApp.CharacterResolveSlugReadModel, error) {
	character, err := uc.r.FindBySlug(cmd.Slug)
	if err != nil {
		return characterApp.CharacterResolveSlugReadModel{}, err
	}
	if character == nil {
		return characterApp.CharacterResolveSlugReadModel{}, ErrCharacterNotFound
	}

	return characterApp.MapDomainToCharacterResolveSlugReadModel(*character), nil
}
