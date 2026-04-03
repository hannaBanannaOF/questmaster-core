package bootstrap

import userTransport "questmaster-core/internal/user/transport/http"

func BuildUserHandler() *userTransport.UserHandler {
	return userTransport.NewUserHandler()
}
