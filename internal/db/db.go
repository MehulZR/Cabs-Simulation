package db

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var DriverStatusKey = "driverStatus:%s"
var PickUpLocationKey = "pickUpLocation:%s"
var DropOffLocationKey = "dropOffLocation:%s"
var DriverCurrentLocationKey = "currentLocation:%s"
var DriverPathKey = "path:%s"

var DriverStatusAvailable = "Available"
var DriverStatusPickingUp = "Picking Up"
var DriverStatusDroppingOff = "Dropping Off"
var DriverStatusWaitingToBeProcessed = "Waiting"

func CreateNewRedisClient() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS_CONN_STRING"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	return rdb
}
