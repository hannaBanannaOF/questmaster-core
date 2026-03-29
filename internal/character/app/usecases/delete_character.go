package character

import (
	characterApp "questmaster-core/internal/character/app"
)

type DeleteCharacterUseCase struct {
	r characterApp.CharacterRepository
}

func NewDeleteCharacter(r characterApp.CharacterRepository) *DeleteCharacterUseCase {
	return &DeleteCharacterUseCase{
		r: r,
	}
}

func (uc *DeleteCharacterUseCase) Execute(cmd characterApp.DeleteCharacterCommand) error {
	character, err := uc.r.FindByID(cmd.ID)
	if err != nil {
		return err
	}

	if character == nil {
		return ErrCharacterNotFound
	}

	if err := character.CanDelete(cmd.UserID); err != nil {
		return err
	}

	deleted, err := uc.r.DeleteByID(character.Id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrCharacterNotFound
	}

	return nil
}
