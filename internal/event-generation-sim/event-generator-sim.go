package eventgenerationsim

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"

	"cabs/internal/queue"
	"cabs/internal/types"
	"cabs/internal/worldMap"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventGenerator struct {
	logger         *slog.Logger
	ctx            context.Context
	rabbitMqClient *queue.RabbitMQClient
	r              *rand.Rand
}

func (e *EventGenerator) Generate() {
	pickUpCoordinates, dropOffCoordinates := types.Coordinate{}, types.Coordinate{}

	for dropOffCoordinates.X == pickUpCoordinates.X && dropOffCoordinates.Y == pickUpCoordinates.Y {
		pickUpCoordinates = worldMap.ValidCoordinates[e.r.IntN(len(worldMap.ValidCoordinates))]
		dropOffCoordinates = worldMap.ValidCoordinates[e.r.IntN(len(worldMap.ValidCoordinates))]
	}

	e.logger.Info(fmt.Sprintf("New ride request: Pickup - %+v DropOff - %+v", pickUpCoordinates, dropOffCoordinates))

	jsonEncodedPayload, err := json.Marshal(types.NewRideRequestEvent{PickUpCoordinates: pickUpCoordinates, DropOffCoordinates: dropOffCoordinates})
	if err != nil {
		e.logger.Error("Failed to marshall data to json")
		panic(err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonEncodedPayload,
	}

	if err = e.rabbitMqClient.Send(e.ctx, queue.NewRideRequestedTopic, msg); err != nil {
		e.logger.Error("Failed to publish NewRideRequestedEvent")
		panic(err)
	}
}

func CreateNewEventGenerator(ctx context.Context, logger *slog.Logger, rabbitMqClient *queue.RabbitMQClient) (*EventGenerator, error) {
	if err := rabbitMqClient.CreateQueue(queue.NewRideRequestedTopic, false, false); err != nil {
		return nil, err
	}

	eventGenerator := EventGenerator{logger, ctx, rabbitMqClient, rand.New(rand.NewPCG(1, 2))}
	return &eventGenerator, nil
}
