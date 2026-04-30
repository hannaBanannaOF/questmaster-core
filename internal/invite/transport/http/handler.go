package invite

import (
	"net/http"
	campaignDomain "questmaster-core/internal/campaign/domain"
	inviteApp "questmaster-core/internal/invite/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

type InviteHandler struct {
	getInviteDetailUC *inviteUsecases.GetInviteDetailUseCase
	acceptInviteUC    *inviteUsecases.AcceptInviteUseCase
	createInviteUC    *inviteUsecases.CreateInviteUseCase
}

func NewInviteHandler(getInviteDetailUC *inviteUsecases.GetInviteDetailUseCase, acceptInviteUC *inviteUsecases.AcceptInviteUseCase, createInviteUC *inviteUsecases.CreateInviteUseCase) *InviteHandler {
	return &InviteHandler{
		getInviteDetailUC: getInviteDetailUC,
		acceptInviteUC:    acceptInviteUC,
		createInviteUC:    createInviteUC,
	}
}

// @Summary Get campaign invite details
// @Description Gets campaign invite details such as available characters, campiagn overview and name, etc
// @Tags v1:invite
// @Param inviteHash path string true "Invite hash"
// @Produce json
// @Success 200 {object} InviteDetailsResponse
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 404 {object} httperrors.HttpError "Invite not found"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/invite/{inviteHash} [get]
func (h *InviteHandler) GetInviteDetails(ctx *context.AppContext) error {
	inviteDetails, err := h.getInviteDetailUC.Execute(inviteApp.GetinviteDetailsCommand{
		UserID: ctx.UserID(),
		Hash:   ctx.InviteHash(),
	})
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, MapInviteDetailsReadModelToResponse(inviteDetails))
	return nil
}

func (h *InviteHandler) AcceptInvite(ctx *context.AppContext) error {
	var body AcceptInviteRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	cmd := MapAcceptRequestToAcceptCommand(body, ctx.UserID(), ctx.InviteHash())

	err := h.acceptInviteUC.Execute(cmd)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusNoContent, nil)
	return nil
}

// @Summary Create campaign invite
// @Description Create a new campaign invite
// @Tags v1:invite
// @Accept json
// @Param request body CreateInviteRequest true "Invite data"
// @Produce json
// @Success 201 {object} InviteCreateResponse
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 403 {object} httperrors.HttpError "Not allowed to create invite"
// @Failure 409 {object} httperrors.HttpError "Invite for campaign already exists"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/invite [post]
func (h *InviteHandler) CreateInvite(ctx *context.AppContext) error {
	var body CreateInviteRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		return httperrors.ErrInvalidRequestBody
	}

	invite, err := h.createInviteUC.Execute(inviteApp.CreateInviteCommand{
		CampaignID: campaignDomain.CampaignID(body.CampaignID),
	})
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusCreated, MapInviteDomainToResponse(invite))
	return nil
}
