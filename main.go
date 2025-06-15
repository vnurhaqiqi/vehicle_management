package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/vnurhaqiqi/vehicle_management/configs"
	"github.com/vnurhaqiqi/vehicle_management/infras"
	"github.com/vnurhaqiqi/vehicle_management/internal/model/dto"
	"github.com/vnurhaqiqi/vehicle_management/internal/repository"
	"github.com/vnurhaqiqi/vehicle_management/internal/service"
	"github.com/vnurhaqiqi/vehicle_management/shared/logger"
)

func main() {
	logger.InitLogger()

	cfg := configs.Get()
	db := infras.ProvidePostgresConn(cfg)

	// RabbitMQ
	if err := infras.ProvideRabbitMQConn(cfg); err != nil {
		log.Error().Err(err).Msg("Error connect to RabbitMQ")
	}
	defer infras.CloseRabbitMQConnection()

	go func() {
		if err := infras.ConsumeMessage(cfg.RabbitMQ.Queues.GeofenceAlert); err != nil {
			log.Error().Err(err).Msg("Error RabbitMQ consumer")
		}
	}()

	vehicleLocationRepository := repository.ProvideVechileLocationRepository(db)
	vehicleLocationService := service.ProvideVechileLocationService(vehicleLocationRepository, cfg)
	// MQTT consumer
	infras.ProvideMQTTClient(cfg, vehicleLocationService.CreateVehicleLocationFromMessage)

	app := fiber.New()

	app.Get("/vehicles/:vehicle_id/location", func(c *fiber.Ctx) error {
		var request dto.ResolveVehicleLocationRequest

		vehicleID := c.Params("vehicle_id")
		if vehicleID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "vechile_id can't be empty",
			})
		}

		request.VehicleID = vehicleID
		resp, err := vehicleLocationService.ResolveByVehicleID(c.Context(), request)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})
		}

		return c.JSON(resp)
	})

	app.Get("/vehicles/:vehicle_id/history", func(c *fiber.Ctx) error {
		var request dto.ResolveVehicleLocationRequest

		vehicleID := c.Params("vehicle_id")
		start := c.Query("start")
		end := c.Query("end")

		if vehicleID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "vechile_id can't be empty",
			})
		}

		request.VehicleID = vehicleID
		request.Start = start
		request.End = end

		resp, err := vehicleLocationService.ResolveByHistory(c.Context(), request)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})
		}

		return c.JSON(resp)
	})

	log.
		Info().
		Msg("Starting server on port " + cfg.Server.Port)

	app.Listen(fmt.Sprintf(":%s", cfg.Server.Port))
}
