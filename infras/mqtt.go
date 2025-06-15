package infras

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/vnurhaqiqi/vehicle_management/configs"
)

func ProvideMQTTClient(
	config *configs.Config,
	handler mqtt.MessageHandler,
) {

	if !config.MQTT.Enabled {
		log.Info().Msg("MQTT Connection is not enabled")
		return
	}

	broker := fmt.Sprintf("tcp://%s:%s", config.MQTT.Host, config.MQTT.Port)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(config.MQTT.ClientID)
	opts.SetDefaultPublishHandler(handler)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("Failed to connect to broker")
		return
	}

	log.Info().Msg("Connected to MQTT broker...")

	if token := client.Subscribe(config.MQTT.Topics.VehicleLocation, 1, nil); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("Failed to subscribe")
		return
	}
}
