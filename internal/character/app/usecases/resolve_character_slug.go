package character

import (
	characterApp "questmaster-core/internal/character/app"
	rpgDomain "questmaster-core/internal/rpg/domain"
)

type ResolveCharacterSlugUseCase struct {
	r characterApp.CharacterRepository
}

func NewResolveCharacterSlug(r characterApp.CharacterRepository) *ResolveCharacterSlugUseCase {
	return &ResolveCharacterSlugUseCase{r: r}
}

func (uc *ResolveCharacterSlugUseCase) Execute(slug rpgDomain.Slug) (characterApp.CharacterResolveSlugReadModel, error) {
	character, err := uc.r.FindBySlug(slug)
	if err != nil {
		return characterApp.CharacterResolveSlugReadModel{}, err
	}
	if character == nil {
		return characterApp.CharacterResolveSlugReadModel{}, ErrCharacterNotFound
	}

	return characterApp.MapDomainToCharacterResolveSlugReadModel(*character), nil
}
