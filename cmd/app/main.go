package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"questmaster-core/cmd/app/bootstrap"
	"questmaster-core/cmd/app/routes"
	rpg "questmaster-core/internal/infra/rpg"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	authHost := os.Getenv("AUTH_HOST")
	realmName := os.Getenv("AUTH_REALM")
	rsa, err := rpg.GetJWKSet(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", authHost, realmName))
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

	// Routes
	router := gin.Default()
	routes.RegisterV1Routes(router, routes.V1RoutesDeps{
		CampaignHandler:  campaignHandler,
		CharacterHandler: characterHandler,
		AuthMiddleware:   rpg.AuthMiddleware(rsa),
		PermMiddleware:   rpg.UpdatePermissionsMiddleware(mongoClient, pgPool),
	})

	// MQTT
	mqttPort, _ := strconv.Atoi(os.Getenv("MQTT_PORT"))

	prod := rpg.RabbitMQProducer{
		Config: rpg.RabbitConfig{
			Host:     os.Getenv("MQTT_HOST"),
			Port:     mqttPort,
			Username: os.Getenv("MQTT_USER"),
			Password: os.Getenv("MQTT_PASSWORD"),
		},
	}
	prod.UpdateGatewayPaths(os.Getenv("GATEWAY_EXCHANGE"), os.Getenv("GATEWAY_URL"))
	runAddr := os.Getenv("RUN_ADDR")
	if runAddr == "" {
		runAddr = "0.0.0.0:8080"
	}
	router.Run(runAddr)
}
