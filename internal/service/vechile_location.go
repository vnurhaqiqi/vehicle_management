package service

import (
	"context"
	"encoding/json"
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/vnurhaqiqi/vehicle_management/configs"
	"github.com/vnurhaqiqi/vehicle_management/infras"
	"github.com/vnurhaqiqi/vehicle_management/internal/model"
	"github.com/vnurhaqiqi/vehicle_management/internal/model/dto"
	"github.com/vnurhaqiqi/vehicle_management/internal/repository"
)

type VechileLocationService interface {
	ResolveByVehicleID(ctx context.Context, request dto.ResolveVehicleLocationRequest) (resp dto.VehicleLocation, err error)
	ResolveByHistory(ctx context.Context, request dto.ResolveVehicleLocationRequest) (resp []dto.VehicleLocation, err error)
	CreateVehicleLocationFromMessage(client mqtt.Client, message mqtt.Message)
}

type VechileLocationServiceImpl struct {
	vehicleLocationRepository repository.VechileLocationRepository
	config                    *configs.Config
}

func ProvideVechileLocationService(
	vehicleLocationRepository repository.VechileLocationRepository,
	config *configs.Config,
) VechileLocationService {
	return &VechileLocationServiceImpl{
		vehicleLocationRepository: vehicleLocationRepository,
		config:                    config,
	}
}

func (s *VechileLocationServiceImpl) ResolveByVehicleID(ctx context.Context, request dto.ResolveVehicleLocationRequest) (resp dto.VehicleLocation, err error) {
	filter := request.ToFilter()
	filter.SetOrderBy("timestamp")
	filter.SetSortBy("DESC")

	vehicleLocations, err := s.vehicleLocationRepository.FindByFilter(ctx, filter)
	if err != nil {
		log.Error().Err(err).Interface("filter", filter).Msg("[ResolveByVehicleID] error vehicleLocationRepository.FindByFilter")
		return
	}

	if len(vehicleLocations) == 0 {
		err = errors.New("not found")
		log.Error().Err(err).Interface("filter", filter).Msg("[ResolveByVehicleID] error vechile locations not found")
		return
	}

	vehicleLocation := vehicleLocations[0]
	resp = dto.NewVehicleLocationResponse(vehicleLocation)

	return
}

func (s *VechileLocationServiceImpl) ResolveByHistory(ctx context.Context, request dto.ResolveVehicleLocationRequest) (resp []dto.VehicleLocation, err error) {
	filter := request.ToFilter()
	filter.SetOrderBy("timestamp")
	filter.SetSortBy("DESC")

	vehicleLocations, err := s.vehicleLocationRepository.FindByFilter(ctx, filter)
	if err != nil {
		log.Error().Err(err).Interface("filter", filter).Msg("[ResolveByVehicleID] error vehicleLocationRepository.FindByFilter")
		return
	}

	if len(vehicleLocations) == 0 {
		err = errors.New("not found")
		log.Error().Err(err).Interface("filter", filter).Msg("[ResolveByVehicleID] error vechile locations not found")
		return
	}

	for _, vechileLocation := range vehicleLocations {
		resp = append(resp, dto.NewVehicleLocationResponse(vechileLocation))
	}

	return
}

func (s *VechileLocationServiceImpl) CreateVehicleLocationFromMessage(client mqtt.Client, message mqtt.Message) {
	log.Info().
		Str("topic", message.Topic()).
		Interface("payload", message.Payload()).
		Msg("[CreateVehicleLocationFromMessage] Received message from topic")

	var payload dto.VehicleLocation

	if err := json.Unmarshal(message.Payload(), &payload); err != nil {
		log.Error().Err(err).Msg("[CreateVehicleLocationFromMessage] Error unmarshaling")
		return
	}

	vehicleLocation := dto.NewVehicleLocationFromRequest(payload)
	err := s.vehicleLocationRepository.Insert(context.Background(), vehicleLocation)
	if err != nil {
		log.Error().Err(err).Interface("payload", payload).Msg("[CreateVehicleLocationFromMessage] error vehicleLocationRepository.Insert")
		return
	}

	log.Info().Str("vehicle_id", payload.VehicleID).Msg("[CreateVehicleLocationFromMessage] Success add new location")

	// publish message to RabbitMQ
	s.sendVehicleLocation(vehicleLocation)
}

func (s *VechileLocationServiceImpl) sendVehicleLocation(vehicleLocation model.VehicleLocation) (err error) {
	message := dto.NewMessageFromVehicleLocation(vehicleLocation)

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Interface("message", message).Msg("[sendVehicleLocation] error marshal json")
		return
	}

	err = infras.PublishMessage(
		s.config.RabbitMQ.Queues.GeofenceAlert,
		s.config.RabbitMQ.Exchange,
		string(jsonMessage),
	)
	if err != nil {
		log.Error().Err(err).Interface("message", jsonMessage).Msg("[sendVehicleLocation] error publish message")
	}

	return
}
