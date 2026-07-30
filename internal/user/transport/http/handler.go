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

// @Summary Get current user profile
// @Description Get current user profile
// @Tags v1:user
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} httperrors.HttpError "Unauthorized - missing or invalid access_token"
// @Failure 500 {object} httperrors.HttpError "Internal server error"
// @Security BearerAuth
// @Router /core/api/v1/user [get]
func (h *UserHandler) GetInfo(ctx *context.AppContext) error {
	ctx.JSON(http.StatusOK, MapUserToUserResponse(ctx.User()))
	return nil
}
