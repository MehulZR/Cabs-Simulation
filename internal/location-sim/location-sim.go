package locationsim

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"cabs/internal/db"
	"cabs/internal/queue"
	"cabs/internal/types"

	wm "cabs/internal/worldMap"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type LocationSimulator struct {
	logger         *slog.Logger
	redis          *redis.Client
	ctx            context.Context
	rabbitMqClient *queue.RabbitMQClient
}

var DRIVERS = []string{"A", "B", "C", "D", "E"}

func (l *LocationSimulator) InitDrivers() {
	// Setting drivers initial location
	for _, driver := range DRIVERS {
		if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverCurrentLocationKey, driver), types.Coordinate{X: 0, Y: 0}).Result(); err != nil {
			l.logger.Error("Failed to set initial location for drivers")
			panic(err)
		}
		if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverStatusKey, driver), types.DriverStatus{Status: db.DriverStatusAvailable}).Result(); err != nil {
			l.logger.Error("Failed to set initial status for drivers")
			panic(err)
		}
	}
}

func (l *LocationSimulator) UpdateDriversLocation() {
	// Updating drivers current location
	for _, driver := range DRIVERS {
		driverStatus := types.DriverStatus{}
		if err := l.redis.HGetAll(l.ctx, fmt.Sprintf(db.DriverStatusKey, driver)).Scan(&driverStatus); err != nil {
			l.logger.Error("Failed to fetch the status of the driver")
			panic(err)
		}

		if driverStatus.Status == db.DriverStatusAvailable {
			continue
		}

		currCoordinates := types.Coordinate{}
		if err := l.redis.HGetAll(l.ctx, fmt.Sprintf(db.DriverCurrentLocationKey, driver)).Scan(&currCoordinates); err != nil {
			l.logger.Error("Failed to fetch current location of the driver")
			panic(err)
		}

		currPath := []types.Coordinate{}

		exists, err := l.redis.HExists(l.ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path").Result()
		if err != nil {
			l.logger.Error("Failed to check whether Path exists for the driver")
			panic(err)
		}

		if exists {
			data, err := l.redis.HGet(l.ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path").Result()
			if err != nil {
				l.logger.Error("Failed to fetch the current path of the driver")
				panic(err)
			}

			if err := json.Unmarshal([]byte(data), &currPath); err != nil {
				l.logger.Error("Failed to unmarshall path of the driver")
				panic(err)
			}
		}

		switch driverStatus.Status {
		case db.DriverStatusPickingUp:
			pickUpCoordinates := types.Coordinate{}
			if err := l.redis.HGetAll(l.ctx, fmt.Sprintf(db.PickUpLocationKey, driver)).Scan(&pickUpCoordinates); err != nil {
				l.logger.Error("Failed to fetch the pickUpLocation of the driver")
				panic(err)
			}

			if len(currPath) == 0 {
				path := calculatePath(currCoordinates, pickUpCoordinates)
				marshalledData, err := json.Marshal(path)
				if err != nil {
					l.logger.Error("Failed to marshall path of the driver")
					panic(err)
				}

				if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path", string(marshalledData)).Result(); err != nil {
					l.logger.Error("Failed to set the path of the driver for pickup")
					panic(err)
				}

				currPath = path
			}

			alreadyThere := len(currPath) == 0

			// nextCoordinate, alreadyThere := calculateNextCoordinate(currCoordinates, pickUpCoordinates)
			if alreadyThere {
				if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverStatusKey, driver), types.DriverStatus{Status: db.DriverStatusDroppingOff}).Result(); err != nil {
					l.logger.Error("Failed to update the status of the driver")
					panic(err)
				}
			} else {
				nextCoordinate, updatedPath := currPath[0], currPath[1:]
				if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverCurrentLocationKey, driver), nextCoordinate).Result(); err != nil {
					l.logger.Error("Failed to update the currentLocation of the driver")
					panic(err)
				}

				marshalledData, err := json.Marshal(updatedPath)
				if err != nil {
					l.logger.Error("Failed to marshall path of the driver")
					panic(err)
				}

				if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path", marshalledData).Result(); err != nil {
					l.logger.Error("Failed to update the path of the driver for pickup")
					panic(err)
				}
			}

		case db.DriverStatusDroppingOff:
			dropOffCoordinates := types.Coordinate{}
			if err := l.redis.HGetAll(l.ctx, fmt.Sprintf(db.DropOffLocationKey, driver)).Scan(&dropOffCoordinates); err != nil {
				l.logger.Error("Failed to fetch the dropOffLocation of the driver")
				panic(err)
			}

			if len(currPath) == 0 {
				path := calculatePath(currCoordinates, dropOffCoordinates)
				marshalledData, err := json.Marshal(path)
				if err != nil {
					l.logger.Error("Failed to marshall path of the driver")
					panic(err)
				}

				if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path", string(marshalledData)).Result(); err != nil {
					l.logger.Error("Failed to set the path of the driver for dropOff")
					panic(err)
				}

				currPath = path
			}

			alreadyThere := len(currPath) == 0
			if !alreadyThere {
				nextCoordinate, updatedPath := currPath[0], currPath[1:]
				if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverCurrentLocationKey, driver), nextCoordinate).Result(); err != nil {
					l.logger.Error("Failed to update the currentLocation of the driver")
					panic(err)
				}

				marshalledData, err := json.Marshal(updatedPath)
				if err != nil {
					l.logger.Error("Failed to marshall path of the driver")
					panic(err)
				}

				if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path", marshalledData).Result(); err != nil {
					l.logger.Error("Failed to update the path of the driver for dropOff")
					panic(err)
				}
			}
		}
	}
}

func calculatePath(current, destination types.Coordinate) []types.Coordinate {
	if current.X == destination.X && current.Y == destination.Y {
		return []types.Coordinate{}
	}

	worldMap := make([][]int, len(wm.WorldMap))
	for i := range wm.WorldMap {
		worldMap[i] = make([]int, len(wm.WorldMap[0]))
		copy(worldMap[i], wm.WorldMap[i])
	}

	type pathDetails struct {
		coordinates types.Coordinate
		path        []types.Coordinate
	}

	directions := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	shortestPath := func(current types.Coordinate) []types.Coordinate {
		worldMap[current.Y][current.X] = 1
		queue := []pathDetails{{coordinates: current, path: []types.Coordinate{}}}

		for len(queue) != 0 {
			var currPathDetails pathDetails
			currPathDetails, queue = queue[0], queue[1:]

			if currPathDetails.coordinates.X == destination.X && currPathDetails.coordinates.Y == destination.Y {
				return currPathDetails.path
			}

			for _, direction := range directions {
				nextX := currPathDetails.coordinates.X + direction[0]
				nextY := currPathDetails.coordinates.Y + direction[1]

				if nextX > 49 || nextX < 0 || nextY > 49 || nextY < 0 || worldMap[nextY][nextX] == 1 {
					continue
				}

				worldMap[nextY][nextX] = 1

				nextCoordinates := types.Coordinate{X: nextX, Y: nextY}
				nextPath := append([]types.Coordinate{}, currPathDetails.path...)
				nextPath = append(nextPath, nextCoordinates)

				queue = append(queue, pathDetails{coordinates: nextCoordinates, path: nextPath})
			}
		}

		return []types.Coordinate{}
	}

	path := shortestPath(current)

	return path
}

func (l *LocationSimulator) TriggerFinishedRides() {
	// Fetch drivers destination
	for _, driver := range DRIVERS {
		driverStatus := types.DriverStatus{}
		if err := l.redis.HGetAll(l.ctx, fmt.Sprintf(db.DriverStatusKey, driver)).Scan(&driverStatus); err != nil {
			l.logger.Error("Failed to fetch the status of the driver")
			panic(err)
		}

		if driverStatus.Status != db.DriverStatusDroppingOff {
			continue
		}

		dropOffLocation := types.Coordinate{}
		if err := l.redis.HGetAll(l.ctx, fmt.Sprintf(db.DropOffLocationKey, driver)).Scan(&dropOffLocation); err != nil {
			l.logger.Error("Failed to fetch dropOffLocation")
			panic(err)
		}

		currentLocation := types.Coordinate{}
		if err := l.redis.HGetAll(l.ctx, fmt.Sprintf(db.DriverCurrentLocationKey, driver)).Scan(&currentLocation); err != nil {
			l.logger.Error("Failed to fetch driver's current location")
			panic(err)
		}

		if dropOffLocation.X == currentLocation.X && dropOffLocation.Y == currentLocation.Y {
			if _, err := l.redis.HSet(l.ctx, fmt.Sprintf(db.DriverStatusKey, driver), types.DriverStatus{Status: db.DriverStatusWaitingToBeProcessed}).Result(); err != nil {
				l.logger.Error("Failed to set the status of the driver")
				panic(err)
			}

			jsonEncodedPayload, err := json.Marshal(types.DriverAvailableEvent{Driver: driver})
			if err != nil {
				l.logger.Error("Failed to marshall json data")
				panic(err)
			}

			msg := amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonEncodedPayload,
			}

			if err = l.rabbitMqClient.Send(l.ctx, queue.RideCompletedTopic, msg); err != nil {
				l.logger.Error("Failed to send RideCompletedEvent to queue")
				panic(err)
			}
		}
	}
}

func CreateNewLocationUpdatorSimulator(ctx context.Context, logger *slog.Logger, redis *redis.Client, rabbitMqClient *queue.RabbitMQClient) (*LocationSimulator, error) {
	if err := rabbitMqClient.CreateQueue(queue.DriverAvailableTopic, false, false); err != nil {
		return nil, err
	}

	locationSimulator := LocationSimulator{logger, redis, ctx, rabbitMqClient}
	return &locationSimulator, nil
}
