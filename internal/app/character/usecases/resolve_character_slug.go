package character

import (
	app "questmaster-core/internal/app/character"
	"questmaster-core/internal/domain/rpg"
)

type ResolveCharacterSlugUseCase struct {
	r app.CharacterRepository
}

func NewResolveCharacterSlug(r app.CharacterRepository) *ResolveCharacterSlugUseCase {
	return &ResolveCharacterSlugUseCase{r: r}
}

func (uc *ResolveCharacterSlugUseCase) Execute(slug string) (int, error) {
	slugDomain, err := rpg.NewSlug(slug)
	if err != nil {
		return 0, err
	}
	domain, err := uc.r.FindBySlug(slugDomain)
	if err != nil {
		return 0, err
	}
	if domain == nil {
		return 0, ErrCharacterNotFound
	}

	return int(domain.Id), nil
}
