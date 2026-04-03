package main

import (
	"context"
	"flag"
	"log"
	"os"
	"questmaster-core/cmd/app/bootstrap"
	"questmaster-core/cmd/app/routes"
	"questmaster-core/internal/shared/middleware"

	_ "questmaster-core/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Debug flags
	debug := flag.Bool("debug", false, "Run in debug mode")
	flag.Parse()
	if *debug {
		log.Print("Running in DEBUG mode, loading envs from file...")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Panicf("Unable to load .env file: %s", err)
		}
	} else {
		log.Print("Running in RELEASE mode!")
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

	mongoClientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))
	mongoClient, err := mongo.Connect(context.Background(), mongoClientOptions)
	if err != nil {
		log.Panicf("error connectiong to MongoDB: %s", err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

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
		middleware.ErrorHandler(),
	)

	routes.RegisterV1Routes(router, routes.V1RoutesDeps{
		CampaignHandler:  campaignHandler,
		CharacterHandler: characterHandler,
		InviteHandler:    inviteHandler,
		UserHandler:      userHandler,
		AuthMiddleware:   middleware.AuthMiddleware(rsa),
		PermMiddleware:   middleware.UpdatePermissionsMiddleware(mongoClient, pgPool),
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	runAddr := os.Getenv("RUN_ADDR")
	if runAddr == "" {
		runAddr = "0.0.0.0:8080"
	}
	router.Run(runAddr)
}
