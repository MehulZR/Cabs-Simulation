package processor

import (
	"cabs/internal/db"
	locationsim "cabs/internal/location-sim"
	"cabs/internal/queue"
	"cabs/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Processor struct {
	ctx      context.Context
	logger   *slog.Logger
	redis    *redis.Client
	rabbitMQ *queue.RabbitMQClient
}

func (p *Processor) Start() {
	if err := p.initializeDriverAvailableQueue(); err != nil {
		p.logger.Error("Failed to initialize driverAvailableQueue")
		panic(err)
	}

	go p.processCompletedRides()

	p.processNewRideRequests()
}

func (p *Processor) processNewRideRequests() {
	newRideRequestedQueue, err := p.rabbitMQ.Consume(queue.NewRideRequestedTopic, "processor-newRideRequested", true)
	if err != nil {
		p.logger.Error("Failed to consume the NewRideRequested queue")
		panic(err)
	}

	driverAvailableQueue, err := p.rabbitMQ.Consume(queue.DriverAvailableTopic, "processor-DriverAvailable", true)
	if err != nil {
		p.logger.Error("Failed to consume the DriverAvailable queue")
		panic(err)
	}

	for driverAvailableEvent := range driverAvailableQueue {
		driverAvailable := types.DriverAvailableEvent{}
		if err = json.Unmarshal(driverAvailableEvent.Body, &driverAvailable); err != nil {
			p.logger.Warn("Failed to unmarshall DriverAvailableEvent")
		}

		newRideRequestedEvent := <-newRideRequestedQueue

		newRideRequest := types.NewRideRequestEvent{}
		if err = json.Unmarshal(newRideRequestedEvent.Body, &newRideRequest); err != nil {
			p.logger.Warn("Failed to unmarshall NewRideRequestEvent")
		}

		if _, err := p.redis.HSet(p.ctx, fmt.Sprintf(db.PickUpLocationKey, driverAvailable.Driver), types.Coordinate{X: newRideRequest.PickUpCoordinates.X, Y: newRideRequest.PickUpCoordinates.Y}).Result(); err != nil {
			p.logger.Error("Failed to set pickUpCoordinates for driver")
			panic(err)
		}

		if _, err := p.redis.HSet(p.ctx, fmt.Sprintf(db.DropOffLocationKey, driverAvailable.Driver), types.Coordinate{X: newRideRequest.DropOffCoordinates.X, Y: newRideRequest.DropOffCoordinates.Y}).Result(); err != nil {
			p.logger.Error("Failed to set dropOffCoordinates for driver")
			panic(err)
		}

		if _, err := p.redis.HSet(p.ctx, fmt.Sprintf(db.DriverStatusKey, driverAvailable.Driver), types.DriverStatus{Status: db.DriverStatusPickingUp}).Result(); err != nil {
			p.logger.Error("Failed to set status for driver")
			panic(err)
		}

		fmt.Printf("Ride request: %+v assigned to driver %s\n", newRideRequest, driverAvailable.Driver)
	}
}

func (p *Processor) initializeDriverAvailableQueue() error {
	for _, driver := range locationsim.DRIVERS {
		jsonEncodedPayload, err := json.Marshal(types.DriverAvailableEvent{Driver: driver})
		if err != nil {
			return err
		}

		msg := amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonEncodedPayload,
		}

		if err := p.rabbitMQ.Send(p.ctx, queue.DriverAvailableTopic, msg); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) processCompletedRides() {
	rabbitMQ, err := queue.CreateNewRabbitMQClient()
	var rabbitMQCloseErr error
	defer func() {
		rabbitMQCloseErr = rabbitMQ.Close()
	}()

	if err != nil {
		p.logger.Error("Failed to create rabbitMqClient")
		panic(err)
	}

	if err := rabbitMQ.CreateQueue(queue.RideCompletedTopic, false, false); err != nil {
		p.logger.Error("Failed to create rideCompletedQueue")
		panic(err)
	}

	RideCompletedQueue, err := rabbitMQ.Consume(queue.RideCompletedTopic, "processor-RideCompleted", true)
	if err != nil {
		p.logger.Error("Failed to consume the RideCompletedQueue")
		panic(err)
	}

	for rideCompletedEvent := range RideCompletedQueue {
		rideCompleted := types.RideCompletedEvent{}
		if err := json.Unmarshal(rideCompletedEvent.Body, &rideCompleted); err != nil {
			p.logger.Error("Failed to unmarshall AvailableDriverEvent")
			panic(err)
		}

		if _, err := p.redis.HSet(p.ctx, fmt.Sprintf(db.DriverStatusKey, rideCompleted.Driver), types.DriverStatus{Status: db.DriverStatusAvailable}).Result(); err != nil {
			p.logger.Error("Failed to reset the driver status for recently freed driver")
			panic(err)
		}

		if _, err := p.redis.HDel(p.ctx, fmt.Sprintf(db.PickUpLocationKey, rideCompleted.Driver), "X", "Y").Result(); err != nil {
			p.logger.Error("Failed to reset pickUpLocation for the recently freed driver")
			panic(err)
		}

		if _, err := p.redis.HDel(p.ctx, fmt.Sprintf(db.DropOffLocationKey, rideCompleted.Driver), "X", "Y").Result(); err != nil {
			p.logger.Error("Failed to reset dropOffLocation for the recently freed driver")
			panic(err)
		}

		jsonEncodedPayload, err := json.Marshal(types.DriverAvailableEvent(rideCompleted))
		if err != nil {
			p.logger.Error("Failed to marshall DriverAvailableEvent")
			panic(err)
		}

		msg := amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonEncodedPayload,
		}

		if err = rabbitMQ.Send(p.ctx, queue.DriverAvailableTopic, msg); err != nil {
			p.logger.Error("Failed to push driverAvailableEvent to kafka")
			panic(err)
		}

		fmt.Printf("Driver got free: %s\n", rideCompleted.Driver)
	}

	if rabbitMQCloseErr != nil {
		p.logger.Error("Failed to close RabbitMQ channel")
	}
}

func CreateNewProcessor(ctx context.Context, logger *slog.Logger, redis *redis.Client, rabbitMQ *queue.RabbitMQClient) (*Processor, error) {
	if err := rabbitMQ.CreateQueue(queue.NewRideRequestedTopic, false, false); err != nil {
		logger.Error("Failed to create NewRideRequested queue")
		panic(err)
	}

	if err := rabbitMQ.CreateQueue(queue.RideCompletedTopic, false, false); err != nil {
		logger.Error("Failed to create RideCompleted queue")
		panic(err)
	}

	if err := rabbitMQ.CreateQueue(queue.DriverAvailableTopic, false, false); err != nil {
		logger.Error("Failed to create DriverAvailable queue")
		panic(err)
	}

	processor := Processor{ctx, logger, redis, rabbitMQ}
	return &processor, nil
}
