package character

import (
	"log"

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

	// Only own characters appear in the character list, so only their opens are recorded.
	// Losing one recency update is harmless; failing the read is not
	if char.IsPlayer(cmd.UserID) {
		if err := uc.r.RecordAccess(cmd.UserID, char.Id); err != nil {
			log.Printf("recording access to character %d: %v", char.Id.Value(), err)
		}
	}

	return *char, nil
}
