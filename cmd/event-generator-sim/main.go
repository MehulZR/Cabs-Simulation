package main

import (
	"context"
	"time"

	eventgenerationsim "cabs/internal/event-generation-sim"

	"cabs/internal/logger"
	"cabs/internal/queue"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()
	log := logger.CreateNewLogger("Event Generator Sim")

	rc, err := queue.CreateNewRabbitMQClient()
	if err != nil {
		log.Error("Failed to create rabbitMqClient")
		panic(err)
	}
	defer rc.Close()
	eventGenerator, err := eventgenerationsim.CreateNewEventGenerator(ctx, log, rc)
	if err != nil {
		log.Error("Failed to create Event Generator Sim")
		panic(err)
	}

	for {
		eventGenerator.Generate()

		time.Sleep(10 * time.Second)
	}
}
