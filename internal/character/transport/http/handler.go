package character

import (
	"net/http"

	characterApp "questmaster-core/internal/character/app"
	characterUsecases "questmaster-core/internal/character/app/usecases"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

type CharactersHandler struct {
	fetchUC           *characterUsecases.FetchMyCharactersUseCase
	createUC          *characterUsecases.CreateCharacterUseCase
	resolveSlugUc     *characterUsecases.ResolveCharacterSlugUseCase
	getDetailsUC      *characterUsecases.GetCharacterDetailUseCase
	updateHPUC        *characterUsecases.UpdateHPUseCase
	deleteCharacterUC *characterUsecases.DeleteCharacterUseCase
}

func NewCharactersHandler(
	fetchUC *characterUsecases.FetchMyCharactersUseCase,
	createUC *characterUsecases.CreateCharacterUseCase,
	resolveSlugUc *characterUsecases.ResolveCharacterSlugUseCase,
	getDetailsUC *characterUsecases.GetCharacterDetailUseCase,
	updateHPUC *characterUsecases.UpdateHPUseCase,
	deleteCharacterUC *characterUsecases.DeleteCharacterUseCase,
) *CharactersHandler {
	return &CharactersHandler{
		fetchUC:           fetchUC,
		createUC:          createUC,
		resolveSlugUc:     resolveSlugUc,
		getDetailsUC:      getDetailsUC,
		updateHPUC:        updateHPUC,
		deleteCharacterUC: deleteCharacterUC,
	}
}

func (h *CharactersHandler) GetMyCharacters(ctx *context.AppContext) error {
	characters, err := h.fetchUC.Execute(ctx.UserID())
	if err != nil {
		return err
	}

	response := make([]CharacterListResponse, 0, len(characters))

	for _, c := range characters {
		response = append(response, MapListReadModelToResponse(c))
	}

	ctx.JSON(http.StatusOK, response)
	return nil
}

func (h *CharactersHandler) CreateCharacter(ctx *context.AppContext) error {
	var body CreateCharacterRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	cmd, err := MapCreateCharacterRequestToCommand(body, ctx.UserID())
	if err != nil {
		return err
	}

	character, err := h.createUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, MapCreateCharacterReadModelToResponse(character))
	return nil
}

func (h *CharactersHandler) ResolveSlug(ctx *context.AppContext) error {
	character, err := h.resolveSlugUc.Execute(ctx.Slug())
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapResolveSlugReadModelToResponse(character))
	return nil
}

func (h *CharactersHandler) GetDetails(ctx *context.AppContext) error {
	details, err := h.getDetailsUC.Execute(ctx.CharacterID())
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapDetailReadModelToResponse(details))
	return nil
}

func (h *CharactersHandler) UpdateHP(ctx *context.AppContext) error {
	var body UpdateHPRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	cmd, err := MapUpdateHPRequestToCommand(body, ctx.CharacterID(), ctx.UserID())
	if err != nil {
		_ = ctx.Error(err)
		return err
	}

	character, err := h.updateHPUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapCurrentHpToResponse(character.CurrentHP))
	return nil
}

func (h *CharactersHandler) DeleteCharacter(ctx *context.AppContext) error {
	cmd := characterApp.DeleteCharacterCommand{
		ID:     ctx.CharacterID(),
		UserID: ctx.UserID(),
	}

	err := h.deleteCharacterUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusNoContent, nil)
	return nil
}
