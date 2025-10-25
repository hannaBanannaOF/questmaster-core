package handlers

import (
	"fmt"
	"log"
	"net/http"
	"questmaster-core/domain/vo"
	"questmaster-core/internal/app/application/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CharacterSheetHandler struct {
	CsSvc services.CharacterSheetServiceInterface
}

func (hdlr *CharacterSheetHandler) GetCharacterSheetDetails(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("characterSheetId"))
	if err != nil {
		log.Panicf("Unable to get character sheet id from params: %s", err)
	}
	result := hdlr.CsSvc.GetCharacterSheetDetails(id)
	if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusNotFound, nil)
	}
}

func (hdlr *CharacterSheetHandler) ResolveSlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	result := hdlr.CsSvc.ResolveSlug(slug)
	if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusNotFound, nil)
	}
}

func (hdlr *CharacterSheetHandler) CreateCharacterSheet(ctx *gin.Context) {
	var body vo.CreateCharacterSheet
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	result := hdlr.CsSvc.CreateCharacterSheet(body.CharacterName, body.MaxHp, body.TrpgSystem, fmt.Sprintf("%s", ctx.Keys["UserID"]))

	ctx.JSON(http.StatusOK, result)
}
