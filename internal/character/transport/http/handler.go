package character

import (
	"net/http"

	characterApp "questmaster-core/internal/character/app"
	characterUsecases "questmaster-core/internal/character/app/usecases"
	characterDomain "questmaster-core/internal/character/domain"
	rpgApp "questmaster-core/internal/rpg/app"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

type CharactersHandler struct {
	getCurrentUserCharactersUC *characterUsecases.GetCurrentUserCharactersUseCase
	createUC                   *characterUsecases.CreateCharacterUseCase
	resolveSlugUc              *characterUsecases.ResolveCharacterSlugUseCase
	getDetailsUC               *characterUsecases.GetCharacterDetailUseCase
	updateHPUC                 *characterUsecases.UpdateHPUseCase
	deleteCharacterUC          *characterUsecases.DeleteCharacterUseCase
	createPremadeUC            *characterUsecases.CreatePremadeCharacterUseCase
}

func NewCharactersHandler(
	getCurrentUserCharactersUC *characterUsecases.GetCurrentUserCharactersUseCase,
	createUC *characterUsecases.CreateCharacterUseCase,
	resolveSlugUc *characterUsecases.ResolveCharacterSlugUseCase,
	getDetailsUC *characterUsecases.GetCharacterDetailUseCase,
	updateHPUC *characterUsecases.UpdateHPUseCase,
	deleteCharacterUC *characterUsecases.DeleteCharacterUseCase,
	createPremadeUC *characterUsecases.CreatePremadeCharacterUseCase,
) *CharactersHandler {
	return &CharactersHandler{
		getCurrentUserCharactersUC: getCurrentUserCharactersUC,
		createUC:                   createUC,
		resolveSlugUc:              resolveSlugUc,
		getDetailsUC:               getDetailsUC,
		updateHPUC:                 updateHPUC,
		deleteCharacterUC:          deleteCharacterUC,
		createPremadeUC:            createPremadeUC,
	}
}

// @Summary Create premade character
// @Description The campaign DM prepares a character in the campaign, with no player, for players to claim when accepting the invite. It gets the campaign game system, and current HP starts at max HP
// @Tags v1:character
// @Param campaignID path integer true "Campaign ID"
// @Accept json
// @Param request body CreatePremadeCharacterRequest true "Premade character data"
// @Produce json
// @Success 201 {object} rpg.RpgSlugResponse
// @Failure 400 {object} httperrors.HttpError "Invalid character data, or the campaign is archived"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 403 {object} httperrors.HttpError "Only the campaign DM can create premade characters"
// @Failure 404 {object} httperrors.HttpError "Campaign not found"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/campaign/{campaignID}/character [post]
func (h *CharactersHandler) CreatePremadeCharacter(ctx *context.AppContext) error {
	var body CreatePremadeCharacterRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	name, err := characterDomain.NewCharacterName(body.Name)
	if err != nil {
		return err
	}

	var hp *characterDomain.HP
	if body.MaxHP != nil {
		h, err := characterDomain.NewHP(*body.MaxHP, *body.MaxHP)
		if err != nil {
			return err
		}
		hp = &h
	}

	character, err := h.createPremadeUC.Execute(characterApp.CreatePremadeCharacterCommand{
		CampaignID: ctx.CampaignID(),
		Name:       name,
		Hp:         hp,
		UserID:     ctx.UserID(),
	})
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, MapCreateCharacterReadModelToResponse(character))
	return nil
}

// @Summary Get current user characters
// @Description Get a page of the current user characters, ordered by `sort` (default: name, ignoring case and accents, then id)
// @Tags v1:character
// @Param game_system query string false "Only characters of this game system"
// @Param without_campaign query boolean false "true: only characters without a campaign; false: only characters in a campaign"
// @Param sort query string false "name: by name; recent: last opened by the user first, never opened after, by name" Enums(name, recent) default(name)
// @Param limit query integer false "Page size, 1 to 100" default(50)
// @Param offset query integer false "Number of characters to skip" default(0)
// @Produce json
// @Success 200 {object} CharacterListPageResponse
// @Failure 400 {object} httperrors.HttpError "Invalid game_system, sort, limit or offset"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/character [get]
func (h *CharactersHandler) GetCurrentUserCharacters(ctx *context.AppContext) error {
	characters, err := h.getCurrentUserCharactersUC.Execute(characterApp.GetCurrentUserCharactersCommand{
		UserID:  ctx.UserID(),
		Filters: ctx.CharacterListFilters(),
		Page:    ctx.Page(),
	})
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapListPageToResponse(characters))
	return nil
}

// @Summary Create character
// @Description Create new character
// @Tags v1:character
// @Accept json
// @Param request body CreateCharacterRequest true "Character data"
// @Produce json
// @Success 201 {object} rpg.RpgSlugResponse
// @Failure 400 {object} httperrors.HttpError "Invalid character data"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/character [post]
func (h *CharactersHandler) CreateCharacter(ctx *context.AppContext) error {

	var body CreateCharacterRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	name, err := characterDomain.NewCharacterName(body.Name)
	if err != nil {
		return err
	}

	system, err := rpgDomain.NewSystem(body.System)
	if err != nil {
		return err
	}

	var hp *characterDomain.HP
	if body.Hp != nil {
		h, err := characterDomain.NewHP(*body.Hp, *body.Hp)
		if err != nil {
			return err
		}
		hp = &h
	}

	cmd := characterApp.CreateCharacterCommand{
		Name:   name,
		System: system,
		Hp:     hp,
		Player: ctx.UserID(),
	}

	character, err := h.createUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, MapCreateCharacterReadModelToResponse(character))
	return nil
}

// @Summary Resolve character slug
// @Description Resolves the character slug to the internal ID
// @Tags v1:character
// @Param slug path string true "Character slug"
// @Produce json
// @Success 200 {object} rpg.RpgIdResponse
// @Failure 400 {object} httperrors.HttpError "Invalid slug"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 404 {object} httperrors.HttpError "Character not found"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/character/resolve/{slug} [get]
func (h *CharactersHandler) ResolveSlug(ctx *context.AppContext) error {
	character, err := h.resolveSlugUc.Execute(rpgApp.ResolveSlugCommand{
		Slug: ctx.Slug(),
	})
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapResolveSlugReadModelToResponse(character))
	return nil
}

// @Summary Get character details
// @Description Get a more detailed view of the character. Visible to its player and, when the character is in a campaign, to the campaign DM and members
// @Tags v1:character
// @Param characterID path integer true "Character ID"
// @Produce json
// @Success 200 {object} CharacterDetailResponse
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 403 {object} httperrors.HttpError "Not the player nor a member of the character campaign"
// @Failure 404 {object} httperrors.HttpError "Character not found"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/character/{characterID} [get]
func (h *CharactersHandler) GetDetails(ctx *context.AppContext) error {
	details, err := h.getDetailsUC.Execute(characterApp.GetCharacterDetailsCommand{
		ID:     ctx.CharacterID(),
		UserID: ctx.UserID(),
	})
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapDetailReadModelToResponse(details, ctx.UserID()))
	return nil
}

// @Summary Update HP
// @Description Updates character HP
// @Tags v1:character
// @Param characterID path integer true "Character ID"
// @Accept json
// @Param request body UpdateHPRequest true "New character HP"
// @Produce json
// @Success 200 {object} CharacterCurrentHpResponse
// @Failure 404 {object} httperrors.HttpError "Character not found"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 403 {object} httperrors.HttpError "Only the character's player can update it"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/character/{characterID}/hp [patch]
func (h *CharactersHandler) UpdateHP(ctx *context.AppContext) error {
	var body UpdateHPRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	cmd, err := MapUpdateHPRequestToCommand(body, ctx.CharacterID(), ctx.UserID())
	if err != nil {
		return err
	}

	character, err := h.updateHPUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapCurrentHpToResponse(character.CurrentHP))
	return nil
}

// @Summary Delete character
// @Description Deletes a character
// @Tags v1:character
// @Param characterID path integer true "Character ID"
// @Produce json
// @Success 204
// @Failure 404 {object} httperrors.HttpError "Character not found"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 403 {object} httperrors.HttpError "Only the player can delete a character"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/character/{characterID} [delete]
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
