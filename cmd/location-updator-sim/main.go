package main

import (
	"context"
	"time"

	"cabs/internal/db"
	locationsim "cabs/internal/location-sim"
	"cabs/internal/logger"
	"cabs/internal/queue"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()
	log := logger.CreateNewLogger("Location Updator Sim")
	redisClient := db.CreateNewRedisClient()

	rc, err := queue.CreateNewRabbitMQClient()
	if err != nil {
		log.Error("Failed to create rabbitMqClient")
		panic(err)
	}
	defer rc.Close()

	locationSimulator, err := locationsim.CreateNewLocationUpdatorSimulator(ctx, log, redisClient, rc)
	if err != nil {
		log.Error("Failed to create Location Simulator")
		panic(err)
	}

	if _, err := redisClient.FlushAll(ctx).Result(); err != nil {
		log.Error("Failed to flush redis db")
		panic(err)
	}

	locationSimulator.InitDrivers()

	for {
		locationSimulator.UpdateDriversLocation()

		locationSimulator.TriggerFinishedRides()

		time.Sleep(1 * time.Second)
	}
}
