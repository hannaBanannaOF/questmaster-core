package handlers

import (
	"fmt"
	"net/http"
	"questmaster-core/internal/app/application/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	SessionSvc        services.SessionServiceInterface
	CharacterSheetSvc services.CharacterSheetServiceInterface
}

func (hdlr *MeHandler) GetMySessions(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, hdlr.SessionSvc.GetAllByPlayerIdOrDmId(fmt.Sprintf("%s", ctx.Keys["UserID"])))
}

func (hdlr *MeHandler) GetMyCalendar(ctx *gin.Context) {
	year, err := strconv.Atoi(ctx.Query("year"))
	if err != nil {
		year = time.Now().Year()
	}
	month, err := strconv.Atoi(ctx.Query("month"))
	if err != nil {
		month = int(time.Now().Month())
	}
	queryDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	ctx.JSON(http.StatusOK, hdlr.SessionSvc.GetCalendar(queryDate, queryDate.AddDate(0, 1, -1), fmt.Sprintf("%s", ctx.Keys["UserID"])))
}

func (hdlr *MeHandler) GetMyUpcoming(ctx *gin.Context) {
	result := hdlr.SessionSvc.GetMyUpcoming(fmt.Sprintf("%s", ctx.Keys["UserID"]))
	if result != nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (hdlr *MeHandler) GetMyCharacterSheets(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, hdlr.CharacterSheetSvc.GetAllByPlayerId(fmt.Sprintf("%s", ctx.Keys["UserID"])))
}
