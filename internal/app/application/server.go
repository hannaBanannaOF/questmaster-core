package server

import (
	"fmt"
	"log"
	"os"
	handlers "questmaster-core/internal/app/application/handlers"
	services_v1 "questmaster-core/internal/app/application/services/v1"
	auth "questmaster-core/internal/app/infra/auth"
	datasource_pg "questmaster-core/internal/app/infra/datasources/impl/postgres"
	db "questmaster-core/internal/app/infra/db"
	mqtt "questmaster-core/internal/app/infra/mqtt"
	"strconv"

	gin "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type QuestmasterServer struct {
	Router *gin.Engine
}

func (svr *QuestmasterServer) Serve() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Unable to load .env file: %s", err)
	}
	authHost := os.Getenv("AUTH_HOST")
	realmName := os.Getenv("AUTH_REALM")
	rsa, err := auth.GetJWKSet(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", authHost, realmName))
	if err != nil {
		log.Panicf("Unable to start server: %s", err)
	}

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	db := db.Db{
		Config: db.DbConfig{
			Host:         os.Getenv("DB_HOST"),
			Port:         dbPort,
			Username:     os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			DatabaseName: os.Getenv("DB_DATABASENAME"),
		}}
	meHdlrV1 := handlers.MeHandler{
		SessionSvc: &services_v1.SessionServiceV1{SessionDs: &datasource_pg.SessionDatasourcePG{
			Db: db,
		}},
		CharacterSheetSvc: &services_v1.CharacterSheetServiceV1{CharacterSheetDs: &datasource_pg.CharacterSheetDataSourcePG{
			Db: db,
		}},
	}
	v1 := svr.Router.Group("/core/api/v1", auth.AuthMiddleware(rsa))
	{
		me := v1.Group("/me")
		{
			me.GET("/sessions", meHdlrV1.GetMySessions)
			me.GET("/sessions/calendar", meHdlrV1.GetMyCalendar)
			me.GET("/sessions/upcoming", meHdlrV1.GetMyUpcoming)
			me.GET("/sheets", meHdlrV1.GetMyCharacterSheets)
		}
	}
	gatewayUrl := os.Getenv("GATEWAY_URL")
	mqttPort, _ := strconv.Atoi(os.Getenv("MQTT_PORT"))

	prod := mqtt.RabbitMQProducer{
		Config: mqtt.RabbitConfig{
			Host:     os.Getenv("MQTT_HOST"),
			Port:     mqttPort,
			Username: os.Getenv("MQTT_USER"),
			Password: os.Getenv("MQTT_PASSWORD"),
		},
	}
	prod.UpdateGatewayPaths("update-paths", gatewayUrl)
	runAddr := os.Getenv("RUN_ADDR")
	svr.Router.Run(runAddr)
}
