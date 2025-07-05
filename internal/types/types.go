package types

type DriverAvailableEvent struct {
	Driver string `json:"driver"`
}

type NewRideRequestEvent struct {
	PickUpCoordinates  Coordinate `json:"pickUp"`
	DropOffCoordinates Coordinate `json:"dropOff"`
}

type RideCompletedEvent struct {
	Driver string `json:"driver"`
}

type Coordinate struct {
	X int `redis:"x"`
	Y int `redis:"y"`
}

type DriverStatus struct {
	Status string `redis:"status"`
}

type DriverState struct {
	Driver             string       `json:"driver"`
	DriverStatus       string       `json:"driverStatus"`
	CurrentCoordinates Coordinate   `json:"currentCoordinates"`
	PickUpCoordinates  Coordinate   `json:"pickUpCoordinates"`
	DropOffCoordinates Coordinate   `json:"dropOffCoordinates"`
	Path               []Coordinate `json:"path"`
}
