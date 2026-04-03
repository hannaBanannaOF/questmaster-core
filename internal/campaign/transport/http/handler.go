package campaign

import (
	"net/http"
	campaignApp "questmaster-core/internal/campaign/app"
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignDomain "questmaster-core/internal/campaign/domain"
	rpg "questmaster-core/internal/rpg/app"
	rpgDomain "questmaster-core/internal/rpg/domain"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

type CampaignHandler struct {
	getCurrentUserCampaignsUC *campaignUsecases.GetCurrentUserCampaignsUseCase
	resolveSlugUC             *campaignUsecases.ResolveCampaignSlugUseCase
	createCampaignUC          *campaignUsecases.CreateCampaignUseCase
	updateStatusUC            *campaignUsecases.UpdateCampaignStatusUseCase
	getDetailsUC              *campaignUsecases.GetCampaignDetailsUseCase
	deleteCampaignUc          *campaignUsecases.DeleteCampaignUseCase
}

func NewCampaignHandler(
	getCurrentUserCampaignsUC *campaignUsecases.GetCurrentUserCampaignsUseCase,
	resolveSlugUC *campaignUsecases.ResolveCampaignSlugUseCase,
	createCampaignUC *campaignUsecases.CreateCampaignUseCase,
	updateStatusUC *campaignUsecases.UpdateCampaignStatusUseCase,
	getDetailsUC *campaignUsecases.GetCampaignDetailsUseCase,
	deleteCampaignUc *campaignUsecases.DeleteCampaignUseCase,
) *CampaignHandler {
	return &CampaignHandler{
		getCurrentUserCampaignsUC: getCurrentUserCampaignsUC,
		resolveSlugUC:             resolveSlugUC,
		createCampaignUC:          createCampaignUC,
		updateStatusUC:            updateStatusUC,
		getDetailsUC:              getDetailsUC,
		deleteCampaignUc:          deleteCampaignUc,
	}
}

// @Summary Get current user campaigns
// @Description Get current user campaigns as **player** and **DM**
// @Tags v1:campaign
// @Produce json
// @Success 200 {object} CampaignListResponse
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/campaign [get]
func (h *CampaignHandler) GetCurrentUserCampaigns(ctx *context.AppContext) error {
	campaigns, err := h.getCurrentUserCampaignsUC.Execute(campaignApp.GetCurrentUserCampaignsCommand{
		UserID: ctx.UserID(),
	})
	if err != nil {
		return err
	}

	response := make([]CampaignListResponse, 0, len(campaigns))

	for _, c := range campaigns {
		response = append(response, MapListReadModelToResponse(c, ctx.UserID()))
	}

	ctx.JSON(http.StatusOK, response)
	return nil
}

// @Summary Resolve campaign slug
// @Description Resolves the campaign slug to the internal ID
// @Tags v1:campaign
// @Param slug path string true "Campaign slug"
// @Produce json
// @Success 200 {object} rpg.RpgIdResponse
// @Failure 400 {object} httperrors.HttpError "Invalid slug"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 404 {object} httperrors.HttpError "Campaign not found"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/campaign/resolve/{slug} [get]
func (h *CampaignHandler) ResolveSlug(ctx *context.AppContext) error {
	campaign, err := h.resolveSlugUC.Execute(rpg.ResolveSlugCommand{
		Slug: ctx.Slug(),
	})
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapResolveSlugReadModelToResponse(campaign))
	return nil
}

// @Summary Create campaign
// @Description Create new campaign with status **DRAFT**
// @Tags v1:campaign
// @Accept json
// @Param request body CreateCampaignRequest true "Campaign data"
// @Produce json
// @Success 201 {object} rpg.RpgSlugResponse
// @Failure 400 {object} httperrors.HttpError "Invalid campaign data"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/campaign [post]
func (h *CampaignHandler) CreateCampaign(ctx *context.AppContext) error {
	var body CreateCampaignRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	name, err := campaignDomain.NewCampaignName(body.Name)
	if err != nil {
		return err
	}

	var overview *campaignDomain.CampaignOverview

	if body.Overview != nil {
		o := campaignDomain.NewCampaignOverview(*body.Overview)
		overview = &o
	}

	system, err := rpgDomain.NewSystem(body.System)
	if err != nil {
		return err
	}

	cmd := campaignApp.CreateCampaignCommand{
		Name:     name,
		Overview: overview,
		DmID:     ctx.UserID(),
		System:   system,
	}

	campaign, err := h.createCampaignUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, MapCreateCampaignReadModelToResponse(campaign))
	return nil
}

// @Summary Update Status
// @Description Updates campaign status to the following rules:
// @Description - **DRAFT** -> ACTIVE
// @Description - **ACTIVE** -> PAUSED / ARCHIVED
// @Description - **PAUSED** -> ACTIVE / ARCHIVED
// @Tags v1:campaign
// @Param campaignID path integer true "Campaign ID"
// @Accept json
// @Param request body UpdateCampaignStatusRequest true "New campaign status"
// @Produce json
// @Success 200 {object} CampaignStatusResponse
// @Failure 404 {object} httperrors.HttpError "Campaign not found"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 403 {object} httperrors.HttpError "Only the DM can update a campaign, or the status transition is invalid"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/campaign/{campaignID}/status [patch]
func (h *CampaignHandler) UpdateStatus(ctx *context.AppContext) error {
	var body UpdateCampaignStatusRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	status, err := campaignDomain.NewCampaignStatus(body.Status)
	if err != nil {
		return err
	}

	cmd := campaignApp.UpdateCampaignStatusCommand{
		CampaignID: ctx.CampaignID(),
		UserID:     ctx.UserID(),
		NewStatus:  status,
	}

	campaign, err := h.updateStatusUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapStatusToResponse(campaign.Status))
	return nil
}

// @Summary Get campaign details
// @Description Get a more detailed view of the campaign
// @Tags v1:campaign
// @Param campaignID path integer true "Campaign ID"
// @Produce json
// @Success 200 {object} CampaignDetailResponse
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 404 {object} httperrors.HttpError "Campaign not found"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/campaign/{campaignID} [get]
func (h *CampaignHandler) GetCampaignDetails(ctx *context.AppContext) error {
	cmd := campaignApp.GetCampaignDetailsCommand{
		ID:     ctx.CampaignID(),
		UserID: ctx.UserID(),
	}

	campaignDetails, err := h.getDetailsUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapDetailReadModelToResponse(campaignDetails))
	return nil
}

// @Summary Delete campaign
// @Description Deletes a campaign with the status **ARCHIVED** or **DRAFT**
// @Tags v1:campaign
// @Param campaignID path integer true "Campaign ID"
// @Produce json
// @Success 204
// @Failure 404 {object} httperrors.HttpError "Campaign not found"
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 403 {object} httperrors.HttpError "Only the DM can delete a campaign"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/campaign/{campaignID} [delete]
func (h *CampaignHandler) DeleteCampaign(ctx *context.AppContext) error {
	cmd := campaignApp.DeleteCampaignCommand{
		ID:     ctx.CampaignID(),
		UserID: ctx.UserID(),
	}

	err := h.deleteCampaignUc.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusNoContent, nil)
	return nil
}
