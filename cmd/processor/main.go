package main

import (
	"cabs/internal/db"
	locationsim "cabs/internal/location-sim"
	"cabs/internal/logger"
	"cabs/internal/processor"
	"cabs/internal/queue"
	"cabs/internal/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	log := logger.CreateNewLogger("Processor")
	redisClient := db.CreateNewRedisClient()
	rabbitMQClient, err := queue.CreateNewRabbitMQClient()
	if err != nil {
		log.Error("Failed to create the RabbitMQ Conn")
		panic(err)
	}
	var rabbitMQClientErr error
	defer func() {
		rabbitMQClientErr = rabbitMQClient.Close()
	}()

	processor, err := processor.CreateNewProcessor(ctx, log, redisClient, rabbitMQClient)
	if err != nil {
		log.Error("Failed to create Processor")
		panic(err)
	}

	go processor.Start()

	serveAPI(ctx, redisClient, rabbitMQClient)

	if rabbitMQClientErr != nil {
		log.Error("Failed to close rabbitMQClient")
		panic(rabbitMQClientErr)
	}
}

var currentSimulationState []types.DriverState

func serveAPI(ctx context.Context, redis *redis.Client, rabbitMQ *queue.RabbitMQClient) {
	if err := rabbitMQ.CreateQueue(queue.NewRideRequestedTopic, false, false); err != nil {
		panic(err)
	}

	go updateCurrentSimulationState(ctx, redis)

	websocketUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws" {
			w.WriteHeader(http.StatusNotFound)
		}
		ws, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		defer ws.Close()

		go func() {
			for {
				if _, _, err := ws.ReadMessage(); err != nil {
					ws.Close()
					break
				}
			}
		}()

		for {
			if err := ws.WriteJSON(currentSimulationState); err != nil {
				if errors.Is(err, websocket.ErrCloseSent) {
					break
				}
			}

			time.Sleep(1 * time.Second)
		}
	})

	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}

func updateCurrentSimulationState(ctx context.Context, redis *redis.Client) {
	for {
		simState := make([]types.DriverState, len(locationsim.DRIVERS))

		for i, driver := range locationsim.DRIVERS {
			driverStatus := types.DriverStatus{}
			if err := redis.HGetAll(ctx, fmt.Sprintf(db.DriverStatusKey, driver)).Scan(&driverStatus); err != nil {
				panic(err)
			}

			pickUpCoordinates := types.Coordinate{}
			if err := redis.HGetAll(ctx, fmt.Sprintf(db.PickUpLocationKey, driver)).Scan(&pickUpCoordinates); err != nil {
				panic(err)
			}

			dropOffCoordinates := types.Coordinate{}
			if err := redis.HGetAll(ctx, fmt.Sprintf(db.DropOffLocationKey, driver)).Scan(&dropOffCoordinates); err != nil {
				panic(err)
			}

			currentCoordinates := types.Coordinate{}
			if err := redis.HGetAll(ctx, fmt.Sprintf(db.DriverCurrentLocationKey, driver)).Scan(&currentCoordinates); err != nil {
				panic(err)
			}

			path := []types.Coordinate{}

			exists, err := redis.HExists(ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path").Result()
			if err != nil {
				panic(err)
			}

			if exists {
				data, err := redis.HGet(ctx, fmt.Sprintf(db.DriverPathKey, driver), "Path").Result()
				if err != nil {
					panic(err)
				}

				if err := json.Unmarshal([]byte(data), &path); err != nil {
					panic(err)
				}
			}

			simState[i] = types.DriverState{Driver: driver, DriverStatus: driverStatus.Status, CurrentCoordinates: currentCoordinates, PickUpCoordinates: pickUpCoordinates, DropOffCoordinates: dropOffCoordinates, Path: path}

		}

		currentSimulationState = simState

		time.Sleep(1 * time.Second)
	}
}
