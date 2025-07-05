build:
	go build -o ./bin/processor ./cmd/processor 
	go build -o ./bin/location-updator-sim ./cmd/location-updator-sim 
	go build -o ./bin/event-generator-sim ./cmd/event-generator-sim 
