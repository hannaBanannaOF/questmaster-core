package character

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	usecases "questmaster-core/internal/app/character/usecases"
	rpg "questmaster-core/internal/transport/rpg/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CharactersHandler struct {
	fetchUC       *usecases.FetchMyCharactersUseCase
	createUC      *usecases.CreateCharacterUseCase
	resolveSlugUc *usecases.ResolveCharacterSlugUseCase
	getDetailsUC  *usecases.GetCharacterDetailUseCase
}

func NewCharactersHandler(
	fetchUC *usecases.FetchMyCharactersUseCase,
	createUC *usecases.CreateCharacterUseCase,
	resolveSlugUc *usecases.ResolveCharacterSlugUseCase,
	getDetailsUC *usecases.GetCharacterDetailUseCase,
) *CharactersHandler {
	return &CharactersHandler{
		fetchUC:       fetchUC,
		createUC:      createUC,
		resolveSlugUc: resolveSlugUc,
		getDetailsUC:  getDetailsUC,
	}
}

func (h *CharactersHandler) GetMyCharacters(ctx *gin.Context) {
	uuid, err := uuid.Parse(ctx.GetString("UserID"))
	if err != nil {
		log.Panicf("Unable to parse userid: %s", err)
	}

	characters, err := h.fetchUC.Execute(uuid)
	if err != nil {
		log.Panicf("Unable to get characters for user: %s", err)
	}

	response := make([]CharacterListResponse, 0, len(characters))

	for _, c := range characters {
		response = append(response, MapListReadModelToResponse(c))
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *CharactersHandler) CreateCharacter(ctx *gin.Context) {
	var body CreateCharacterRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := uuid.Parse(ctx.GetString("UserID"))
	if err != nil {
		log.Panicf("Unable to parse userid: %s", err)
	}

	slug, err := h.createUC.Execute(body.Name, body.System, body.Hp, userId)
	if err != nil {
		log.Panicf("Unable to create character: %s", err)
	}

	ctx.JSON(http.StatusOK, rpg.MapSlugToResponse(slug))
}

func (h *CharactersHandler) ResolveSlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	campaignId, err := h.resolveSlugUc.Execute(slug)
	if err != nil {
		if errors.Is(err, usecases.ErrCharacterNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Panicf("Unable to resolve character slug: %s", err)
	}

	ctx.JSON(http.StatusOK, rpg.MapIdToResponse(campaignId))
}

func (h *CharactersHandler) GetDetails(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("characterId"))
	if err != nil {
		log.Panicf("Unable to parse character id: %s", err)
	}

	details, err := h.getDetailsUC.Execute(id)
	if err != nil {
		if errors.Is(err, usecases.ErrCharacterNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Panicf("Unable to get character details: %s", err)
	}

	ctx.JSON(http.StatusOK, MapDetailReadModelToResponse(details))
}
