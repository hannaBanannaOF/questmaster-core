package mqtt

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	Config RabbitConfig
}

type RabbitConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type UpdatePathsMessage struct {
	Instance string `json:"instance"`
	Url      string `json:"url"`
	Regex    string `json:"regex"`
}

func (prod *RabbitMQProducer) UpdateGatewayPaths(ExchangeName string, GatewayUrl string) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", prod.Config.Username, prod.Config.Password, prod.Config.Host, prod.Config.Port))

	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
	}
	if conn != nil {
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel: %s", err)
		}
		defer ch.Close()

		err = ch.ExchangeDeclare(ExchangeName, "direct", true, false, false, false, nil)
		if err != nil {
			log.Printf("Failed to declare exchange: %s", err)
		}

		data := UpdatePathsMessage{
			Instance: "questmaster-core",
			Url:      GatewayUrl,
			Regex:    "/core/**",
		}

		body, err := json.Marshal(data)
		if err != nil {
			log.Printf("Failed to marshal JSON: %s", err)
		}

		err = ch.Publish("update-paths", "questmaster-core", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		if err != nil {
			log.Printf("Failed to publish a message: %s", err)
		}
		log.Printf("Gateway paths updated!")
	}

}
