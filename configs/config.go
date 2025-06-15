package configs

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"PORT"`
	} `mapstructure:"SERVER"`
	DB struct {
		Postgres struct {
			Host            string        `mapstructure:"HOST"`
			Port            string        `mapstructure:"PORT"`
			User            string        `mapstructure:"USER"`
			Password        string        `mapstructure:"PASSWORD"`
			Name            string        `mapstructure:"NAME"`
			MaxConnLifetime time.Duration `mapstructure:"MAX_CONNECTION_LIFETIME"`
			MaxIdleConn     int           `mapstructure:"MAX_IDLE_CONNECTION"`
			MaxOpenConn     int           `mapstructure:"MAX_OPEN_CONNECTION"`
		} `mapstructure:"PG"`
	} `mapstructure:"DB"`
	MQTT struct {
		Enabled  bool   `mapstructure:"ENABLED"`
		Host     string `mapstructure:"HOST"`
		Port     string `mapstructure:"PORT"`
		ClientID string `mapstructure:"CLIENT_ID"`
		Topics   struct {
			VehicleLocation string `mapstructure:"VEHICLE_LOCATION"`
		} `mapstructure:"TOPICS"`
	} `mapstructure:"MQTT"`
	RabbitMQ struct {
		Enabled  bool   `mapstructure:"ENABLED"`
		Host     string `mapstructure:"HOST"`
		Port     string `mapstructure:"PORT"`
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
		Queues   struct {
			GeofenceAlert string `mapstructure:"GEOFENCE_ALERT"`
		} `mapstructure:"QUEUES"`
		Exchange string `mapstructure:"EXCHANGE"`
	} `mapstructure:"RABBITMQ"`
}

var (
	conf Config
	once sync.Once
)

func Get() *Config {

	once.Do(func() {
		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed reading config file")
		}

		log.Info().
			Msg("Service configuration initialized.")

		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("")
		}
	})

	return &conf
}
