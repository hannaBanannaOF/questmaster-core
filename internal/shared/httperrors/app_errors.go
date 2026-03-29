package httperrors

import (
	"errors"
	"fmt"
	"net/http"
	campaignAppErr "questmaster-core/internal/campaign/app/usecases"
	campaignDomainErr "questmaster-core/internal/campaign/domain"
	characterAppErr "questmaster-core/internal/character/app/usecases"
	characterDomainErr "questmaster-core/internal/character/domain"
	inviteAppErr "questmaster-core/internal/invite/app/usecases"
)

var ErrInvalidParam = errors.New("Invalid param")
var ErrUnauthorized = errors.New("Unauthorized")
var ErrInvalidRequestBody = errors.New("Invalid request body")

type HttpError struct {
	Status  int
	Message string
}

func From(err error) HttpError {
	switch {
	case errors.Is(err, campaignAppErr.ErrCampaignNotFound),
		errors.Is(err, characterAppErr.ErrCharacterNotFound),
		errors.Is(err, inviteAppErr.ErrInviteNotFound):
		return HttpError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		}
	case errors.Is(err, campaignDomainErr.ErrNotDM),
		errors.Is(err, campaignDomainErr.ErrNotDeletableStatus),
		errors.Is(err, campaignDomainErr.ErrInvalidStatusTransition),
		errors.Is(err, characterDomainErr.ErrNotPlayer),
		errors.Is(err, characterDomainErr.ErrNotAllowed):
		return HttpError{
			Status:  http.StatusForbidden,
			Message: err.Error(),
		}
	case errors.Is(err, characterDomainErr.ErrInvalidCurrentHP),
		errors.Is(err, campaignDomainErr.ErrEmptyCampaignName),
		errors.Is(err, campaignDomainErr.ErrInvalidCampaignStatus),
		errors.Is(err, characterDomainErr.ErrInvalidMaxHP),
		errors.Is(err, characterDomainErr.ErrInvalidCharacterName),
		errors.Is(err, ErrInvalidParam),
		errors.Is(err, ErrInvalidRequestBody):
		return HttpError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	case errors.Is(err, inviteAppErr.ErrInviteAlreadyExists),
		errors.Is(err, characterAppErr.ErrAlreadyEnrolled):
		return HttpError{
			Status:  http.StatusConflict,
			Message: err.Error(),
		}
	case errors.Is(err, ErrUnauthorized):
		return HttpError{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		}
	default:
		return HttpError{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Internal server error: %s", err),
		}
	}
}
