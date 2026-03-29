package invite

import (
	"net/http"
	inviteApp "questmaster-core/internal/invite/app"
	inviteUsecases "questmaster-core/internal/invite/app/usecases"
	"questmaster-core/internal/shared/context"
	"questmaster-core/internal/shared/httperrors"
)

type InviteHandler struct {
	getInviteDetailUC *inviteUsecases.GetInviteDetailUseCase
	acceptInviteUC    *inviteUsecases.AcceptInviteUseCase
}

func NewInviteHandler(getInviteDetailUC *inviteUsecases.GetInviteDetailUseCase, acceptInviteUC *inviteUsecases.AcceptInviteUseCase) *InviteHandler {
	return &InviteHandler{
		getInviteDetailUC: getInviteDetailUC,
		acceptInviteUC:    acceptInviteUC,
	}
}

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
