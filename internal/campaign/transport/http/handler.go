package campaign

import (
	"errors"
	"net/http"
	campaignApp "questmaster-core/internal/campaign/app"
	campaignUsecases "questmaster-core/internal/campaign/app/usecases"
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

type CampaignHandler struct {
	fetchUC                     *campaignUsecases.FetchMyCampaignsUseCase
	resolveSlugUC               *campaignUsecases.ResolveCampaignSlugUseCase
	createCampaignUC            *campaignUsecases.CreateCampaignUseCase
	updateStatusUC              *campaignUsecases.UpdateCampaignStatusUseCase
	getDetailsUC                *campaignUsecases.GetCampaignDetailsUseCase
	deleteCampaignUc            *campaignUsecases.DeleteCampaignUseCase
	getOrCreateCampaignInviteUc *campaignUsecases.GetOrCreateCampaignInviteUseCase
}

func NewCampaignHandler(
	fetchUC *campaignUsecases.FetchMyCampaignsUseCase,
	resolveSlugUC *campaignUsecases.ResolveCampaignSlugUseCase,
	createCampaignUC *campaignUsecases.CreateCampaignUseCase,
	updateStatusUC *campaignUsecases.UpdateCampaignStatusUseCase,
	getDetailsUC *campaignUsecases.GetCampaignDetailsUseCase,
	deleteCampaignUc *campaignUsecases.DeleteCampaignUseCase,
	getOrCreateCampaignInviteUc *campaignUsecases.GetOrCreateCampaignInviteUseCase,
) *CampaignHandler {
	return &CampaignHandler{
		fetchUC:                     fetchUC,
		resolveSlugUC:               resolveSlugUC,
		createCampaignUC:            createCampaignUC,
		updateStatusUC:              updateStatusUC,
		getDetailsUC:                getDetailsUC,
		deleteCampaignUc:            deleteCampaignUc,
		getOrCreateCampaignInviteUc: getOrCreateCampaignInviteUc,
	}
}

func (h *CampaignHandler) GetMyCampaigns(ctx *context.AppContext) error {
	campaigns, err := h.fetchUC.Execute(ctx.UserID())
	if err != nil {
		return err
	}

	response := make([]CampaignListResponse, 0, len(campaigns))

	for _, c := range campaigns {
		response = append(response, MapListReadModelToResponse(c))
	}

	ctx.JSON(http.StatusOK, response)
	return nil
}

func (h *CampaignHandler) ResolveSlug(ctx *context.AppContext) error {
	campaign, err := h.resolveSlugUC.Execute(ctx.Slug())
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapResolveSlugReadModelToResponse(campaign))
	return nil
}

func (h *CampaignHandler) CreateCampaign(ctx *context.AppContext) error {
	var body CreateCampaignRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	cmd, err := MapCreateRequestToCreateCommand(body, ctx.UserID())
	if err != nil {
		return err
	}

	campaign, err := h.createCampaignUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, MapCreateCampaignReadModelToResponse(campaign))
	return nil
}

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

func (h *CampaignHandler) GetOrCreateCampaignInvite(ctx *context.AppContext) error {
	cmd := campaignApp.GetOrCreateCampaignInviteCommand{
		CampaignID: ctx.CampaignID(),
		UserID:     ctx.UserID(),
	}

	invite, err := h.getOrCreateCampaignInviteUc.Execute(cmd)
	if err != nil {

		if errors.Is(err, inviteUsecases.ErrInviteAlreadyExists) {
			ctx.JSON(http.StatusOK, MapGetOrCreateCampaignInviteReadModelToResponse(invite))
			return nil
		}

		return err
	}

	ctx.JSON(http.StatusCreated, MapGetOrCreateCampaignInviteReadModelToResponse(invite))
	return nil
}
