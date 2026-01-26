package campaign

import (
	"errors"
	"log"
	"net/http"
	usecases "questmaster-core/internal/app/campaign/usecases"
	domain "questmaster-core/internal/domain/campaign"
	rpg "questmaster-core/internal/transport/rpg/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CampaignsHandler struct {
	fetchUC          *usecases.FetchMyCampaignsUseCase
	resolveSlugUC    *usecases.ResolveCampaignSlugUseCase
	createCampaignUC *usecases.CreateCampaignUseCase
	updateStatusUC   *usecases.UpdateCampaignStatusUseCase
}

func NewCampaignsHandler(
	fetchUC *usecases.FetchMyCampaignsUseCase,
	resolveSlugUC *usecases.ResolveCampaignSlugUseCase,
	createCampaignUC *usecases.CreateCampaignUseCase,
	updateStatusUC *usecases.UpdateCampaignStatusUseCase,
) *CampaignsHandler {
	return &CampaignsHandler{
		fetchUC:          fetchUC,
		resolveSlugUC:    resolveSlugUC,
		createCampaignUC: createCampaignUC,
		updateStatusUC:   updateStatusUC,
	}
}

func (h *CampaignsHandler) GetMyCampaigns(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.GetString("UserID"))
	if err != nil {
		log.Panicf("Unable to parse userid: %s", err)
	}

	campaigns, err := h.fetchUC.Execute(userId)
	if err != nil {
		log.Panicf("Unable to get campaigns for user: %s", err)
	}

	response := make([]CampaignListResponse, 0, len(campaigns))

	for _, c := range campaigns {
		response = append(response, MapListReadModelToResponse(c))
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *CampaignsHandler) ResolveSlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	campaignId, err := h.resolveSlugUC.Execute(slug)
	if err != nil {
		if errors.Is(err, usecases.ErrCampaignNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Panicf("Unable to resolve campaign slug: %s", err)
	}

	ctx.JSON(http.StatusOK, rpg.MapIdToResponse(campaignId))
}

func (h *CampaignsHandler) CreateCampaign(ctx *gin.Context) {
	var body CreateCampaignRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := uuid.Parse(ctx.GetString("UserID"))
	if err != nil {
		log.Panicf("Unable to parse userid: %s", err)
	}

	slug, err := h.createCampaignUC.Execute(body.Name, body.Overview, body.System, userId)
	if err != nil {
		log.Panicf("Unable to create campaign: %s", err)
	}

	ctx.JSON(http.StatusOK, rpg.MapSlugToResponse(slug))
}

func (h *CampaignsHandler) UpdateStatus(ctx *gin.Context) {
	var body UpdateCampaignStatusRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	campaignId, err := strconv.Atoi(ctx.Param("campaignId"))
	if err != nil {
		log.Panicf("Unable to parse campaignId: %s", err)
	}

	userId, err := uuid.Parse(ctx.GetString("UserID"))
	if err != nil {
		log.Panicf("Unable to parse userid: %s", err)
	}

	newStatus, err := h.updateStatusUC.Execute(campaignId, userId, body.Status)
	if err != nil {
		if errors.Is(err, domain.ErrNotDM) || errors.Is(err, domain.ErrInvalidStatusTransition) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, usecases.ErrCampaignNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Panicf("UInable to update campaign status: %s", err)
	}

	ctx.JSON(http.StatusOK, MapStatusToResponse(newStatus))

}
