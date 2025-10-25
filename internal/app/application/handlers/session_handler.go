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

type SessionHandler struct {
	SessionSvc services.SessionServiceInterface
}

func (hdlr *SessionHandler) GetSessionDetails(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("sessionId"))
	if err != nil {
		log.Panicf("Unable to get session id from params: %s", err)
	}
	result := hdlr.SessionSvc.GetSessionDetails(id, fmt.Sprintf("%s", ctx.Keys["UserID"]))
	if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusNotFound, nil)
	}
}

func (hdlr *SessionHandler) ToggleSessionInPlay(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("sessionId"))
	if err != nil {
		log.Panicf("Unable to get session id from params: %s", err)
	}
	result := hdlr.SessionSvc.ToggleInPlayById(id, fmt.Sprintf("%s", ctx.Keys["UserID"]))
	if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (hdlr *SessionHandler) ResolveSlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	result := hdlr.SessionSvc.ResolveSlug(slug)
	if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusNotFound, nil)
	}
}

func (hdlr *SessionHandler) CreateSession(ctx *gin.Context) {
	var body vo.CreateSession
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	result := hdlr.SessionSvc.CreateSession(body.SessionName, body.SessionOverview, body.TrpgSystem, fmt.Sprintf("%s", ctx.Keys["UserID"]))

	ctx.JSON(http.StatusOK, result)
}
