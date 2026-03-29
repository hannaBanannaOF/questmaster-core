package session

import "questmaster-core/internal/shared/context"

type SessionHandler struct {
}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{}
}

func (h *SessionHandler) GetUpcomingSession(ctx context.AppContext) error {
	return nil
}
