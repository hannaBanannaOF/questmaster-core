package user

import (
	"net/http"
	"questmaster-core/internal/shared/context"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetInfo(ctx *context.AppContext) error {
	ctx.JSON(http.StatusOK, MapUserToUserResponse(ctx.User()))
	return nil
}
