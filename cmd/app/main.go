package main

import (
	server "questmaster-core/internal/app/application"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	server := server.QuestmasterServer{Router: router}
	server.Serve()
}
