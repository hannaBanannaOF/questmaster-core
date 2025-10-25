package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	handlers "questmaster-core/internal/app/application/handlers"
	services_v1 "questmaster-core/internal/app/application/services/v1"
	auth "questmaster-core/internal/app/infra/auth"
	datasource_pg "questmaster-core/internal/app/infra/datasources/impl/postgres"
	mqtt "questmaster-core/internal/app/infra/mqtt"
	"questmaster-core/internal/app/infra/perm"
	"strconv"

	gin "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuestmasterServer struct {
	Router *gin.Engine
}

func (svr *QuestmasterServer) Serve() {
	debug := flag.Bool("debug", false, "Run in debug mode")
	flag.Parse()
	if *debug {
		log.Print("Running in DEBUG mode, loading envs from file...")
		err := godotenv.Load(".env")
		if err != nil {
			log.Panicf("Unable to load .env file: %s", err)
		}
	} else {
		log.Print("Running in RELEASE mode!")
	}
	authHost := os.Getenv("AUTH_HOST")
	realmName := os.Getenv("AUTH_REALM")
	rsa, err := auth.GetJWKSet(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", authHost, realmName))
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

	// Datasource
	characterSheetDs := datasource_pg.CharacterSheetDataSourcePG{
		Db: pgPool,
	}

	//Services
	sessionSvc := services_v1.SessionServiceV1{
		SessionDs: &datasource_pg.SessionDatasourcePG{
			Db: pgPool,
		},
		CharacterSheetDs: &characterSheetDs,
	}
	characterSheetSvc := services_v1.CharacterSheetServiceV1{
		CharacterSheetDs: &characterSheetDs,
	}

	//Handlers
	meHdlrV1 := handlers.MeHandler{
		SessionSvc:        &sessionSvc,
		CharacterSheetSvc: &characterSheetSvc,
	}
	sessionHdlrV1 := handlers.SessionHandler{
		SessionSvc: &sessionSvc,
	}
	characterSheetHdlrV1 := handlers.CharacterSheetHandler{
		CsSvc: &characterSheetSvc,
	}

	//Routes
	v1 := svr.Router.Group("/core/api/v1", auth.AuthMiddleware(rsa), perm.UpdatePermissionsMiddleware(mongoClient, pgPool))
	{
		me := v1.Group("/me")
		{
			me.GET("/sessions", meHdlrV1.GetMySessions)
			me.GET("/sessions/calendar", meHdlrV1.GetMyCalendar)
			me.GET("/sessions/upcoming", meHdlrV1.GetMyUpcoming)
			me.GET("/sheets", meHdlrV1.GetMyCharacterSheets)
		}
		session := v1.Group("/session")
		{
			session.POST("", sessionHdlrV1.CreateSession)
			session.GET("/resolve/:slug", sessionHdlrV1.ResolveSlug)
			session.GET("/:sessionId", perm.CheckViewSessionPermission(pgPool), sessionHdlrV1.GetSessionDetails)
			session.PUT("/:sessionId/toggle-in-play", perm.CheckUpdateSessionPermission(pgPool), sessionHdlrV1.ToggleSessionInPlay)
		}
		characterSheet := v1.Group("/character-sheet")
		{
			characterSheet.POST("", characterSheetHdlrV1.CreateCharacterSheet)
			characterSheet.GET("/resolve/:slug", characterSheetHdlrV1.ResolveSlug)
			characterSheet.GET("/:characterSheetId", perm.CheckViewCharacterSheetPermission(pgPool), characterSheetHdlrV1.GetCharacterSheetDetails)
		}
	}
	mqttPort, _ := strconv.Atoi(os.Getenv("MQTT_PORT"))

	prod := mqtt.RabbitMQProducer{
		Config: mqtt.RabbitConfig{
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
	svr.Router.Run(runAddr)
}
