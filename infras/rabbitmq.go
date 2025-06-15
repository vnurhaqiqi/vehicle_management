package infras

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/vnurhaqiqi/vehicle_management/configs"
)

var (
	Channel    *amqp.Channel
	Connection *amqp.Connection
)

func ProvideRabbitMQConn(config *configs.Config) error {
	if !config.RabbitMQ.Enabled {
		log.Info().Msg("RabbitMQ is not enabled")
		return nil
	}

	uri := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
	)

	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Error().Err(err).Msg("[RabbitMQ] Failed to connect RabbitMQ")
		return err
	}

	log.Info().Msg("Connected to RabbitMQ...")

	Connection = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("{RabbitMQ] Failed to open a channel")
		return err
	}

	Channel = ch

	return nil
}

func DeclareQueue(queueName string) (amqp.Queue, error) {
	queue, err := Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Err(err).Msg("[RabbitMQ] Failed to declare queue")
		return amqp.Queue{}, err
	}

	return queue, nil
}

func PublishMessage(queueName, exchange, message string) error {
	queue, err := DeclareQueue(queueName)
	if err != nil {
		return err
	}

	err = Channel.Publish(
		exchange,
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		log.Error().Err(err).Msg("[RabbitMQ] Failed to publish message")
		return err
	}

	log.Info().Interface("message", message).Msg("[RabbitMQ] Sent message")
	fmt.Println(message)

	return nil
}

func ConsumeMessage(queueName string) error {
	_, err := DeclareQueue(queueName)
	if err != nil {
		return err
	}

	messagesFromTopic, err := Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error().Err(err).Msg("[RabbitMQ] Failed to register consumer")
		return err
	}

	forever := make(chan bool)

	go func() {
		for m := range messagesFromTopic {
			log.Info().Str("message", string(m.Body)).Msg("[RabbitMQ] Receive Message")
			fmt.Println(string(m.Body))
		}
	}()

	log.Info().Str("queueName", queueName).Msg("[RabbitMQ] Waiting for message on queue")

	<-forever
	return nil

}

func CloseRabbitMQConnection() {
	if Channel != nil {
		Channel.Close()
	}
	if Connection != nil {
		Connection.Close()
	}
}
