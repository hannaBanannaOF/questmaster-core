package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"questmaster-core/cmd/app/bootstrap"
	"questmaster-core/cmd/app/routes"
	"questmaster-core/internal/shared/middleware"

	_ "questmaster-core/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TO RUN SWAGGER -> go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/app/main.go --parseInternal

// @title Questmaster's APIs
// @version 1.0
// @description REST APIs of Questmaster for documentation.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8082
// @BasePath /core/api/v1
// @schemes http
func main() {
	runAddr := os.Getenv("RUN_ADDR")
	if runAddr == "" {
		runAddr = "0.0.0.0:8080"
	}
	healthCheck := flag.Bool("check-health", false, "Executa o healthcheck")
	flag.Parse()

	if *healthCheck {
		_, port, _ := net.SplitHostPort(runAddr)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	oidcHost := os.Getenv("OIDC_HOST")
	rsa, err := middleware.GetJWKSet(oidcHost)
	if err != nil {
		log.Panicf("Unable to start server: %s", err)
	}

	// Database connections
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	if err != nil {
		log.Panicf("unable to get conn pool: %s", err)
	}
	pgPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Panicf("unable to get conn pool: %s", err)
	}
	defer pgPool.Close()

	// Bootstrap
	campaignHandler := bootstrap.BuildCampaignHandler(pgPool)
	characterHandler := bootstrap.BuildCharacterHandler(pgPool)
	inviteHandler := bootstrap.BuildInviteHandler(pgPool)
	userHandler := bootstrap.BuildUserHandler()

	// Routes
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.ErrorHandlerMiddleware(),
		middleware.QueryParamsMiddleware(),
	)

	routes.RegisterV1Routes(router, routes.V1RoutesDeps{
		CampaignHandler:  campaignHandler,
		CharacterHandler: characterHandler,
		InviteHandler:    inviteHandler,
		UserHandler:      userHandler,
		AuthMiddleware:   middleware.AuthMiddleware(rsa),
	})

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(runAddr)
}
