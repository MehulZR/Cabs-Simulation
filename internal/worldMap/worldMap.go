package worldMap

import (
	"cabs/internal/types"
)

var ValidCoordinates []types.Coordinate
var WorldMap [][]int

func init() {
	WorldMap = make([][]int, 50)
	for i := range 50 {
		WorldMap[i] = make([]int, 50)
	}

	for _, obstacle := range Obstacles {
		for y := obstacle.yStart; y < obstacle.yEnd; y++ {
			for x := obstacle.xStart; x < obstacle.xEnd; x++ {
				WorldMap[y][x] = 1
			}
		}
	}

	for y := range 50 {
		for x := range 50 {
			if WorldMap[y][x] != 1 {
				ValidCoordinates = append(ValidCoordinates, types.Coordinate{X: x, Y: y})
			}
		}
	}
}

type ObstacleType string

const (
	Building ObstacleType = "building"
	Park     ObstacleType = "park"
	River    ObstacleType = "river"
)

type Obstacle struct {
	xStart       int
	xEnd         int
	yStart       int
	yEnd         int
	obstacleType ObstacleType
}

var Obstacles = []Obstacle{
	{xStart: 0, xEnd: 5, yStart: 1, yEnd: 7, obstacleType: Building},
	{xStart: 40, xEnd: 49, yStart: 17, yEnd: 23, obstacleType: Park},
	{xStart: 6, xEnd: 9, yStart: 0, yEnd: 5, obstacleType: Building},
	{xStart: 6, xEnd: 9, yStart: 6, yEnd: 10, obstacleType: Building},
	{xStart: 10, xEnd: 16, yStart: 0, yEnd: 10, obstacleType: Building},
	{xStart: 17, xEnd: 18, yStart: 0, yEnd: 2, obstacleType: Building},
	{xStart: 17, xEnd: 18, yStart: 3, yEnd: 7, obstacleType: Building},
	{xStart: 17, xEnd: 18, yStart: 8, yEnd: 12, obstacleType: Building},
	{xStart: 19, xEnd: 25, yStart: 1, yEnd: 12, obstacleType: Park},
	{xStart: 26, xEnd: 29, yStart: 0, yEnd: 13, obstacleType: River},
	{xStart: 26, xEnd: 29, yStart: 14, yEnd: 15, obstacleType: River},
	{xStart: 26, xEnd: 29, yStart: 16, yEnd: 28, obstacleType: River},
	{xStart: 0, xEnd: 15, yStart: 28, yEnd: 30, obstacleType: River},
	{xStart: 16, xEnd: 17, yStart: 28, yEnd: 30, obstacleType: River},
	{xStart: 18, xEnd: 40, yStart: 28, yEnd: 30, obstacleType: River},
	{xStart: 39, xEnd: 41, yStart: 28, yEnd: 50, obstacleType: River},
	{xStart: 0, xEnd: 5, yStart: 8, yEnd: 10, obstacleType: Park},
	{xStart: 0, xEnd: 16, yStart: 11, yEnd: 15, obstacleType: Building},
	{xStart: 16, xEnd: 17, yStart: 13, yEnd: 15, obstacleType: Building},
	{xStart: 18, xEnd: 25, yStart: 13, yEnd: 15, obstacleType: Building},
	{xStart: 0, xEnd: 14, yStart: 16, yEnd: 27, obstacleType: Park},
	{xStart: 15, xEnd: 22, yStart: 16, yEnd: 27, obstacleType: Building},
	{xStart: 23, xEnd: 25, yStart: 16, yEnd: 27, obstacleType: Building},
	{xStart: 30, xEnd: 50, yStart: 0, yEnd: 2, obstacleType: Building},
	{xStart: 30, xEnd: 39, yStart: 3, yEnd: 6, obstacleType: Building},
	{xStart: 40, xEnd: 50, yStart: 3, yEnd: 6, obstacleType: Building},
	{xStart: 30, xEnd: 34, yStart: 7, yEnd: 10, obstacleType: Building},
	{xStart: 35, xEnd: 43, yStart: 7, yEnd: 10, obstacleType: Park},
	{xStart: 44, xEnd: 50, yStart: 7, yEnd: 10, obstacleType: Building},
	{xStart: 30, xEnd: 48, yStart: 11, yEnd: 16, obstacleType: Building},
	{xStart: 49, xEnd: 50, yStart: 11, yEnd: 16, obstacleType: Building},
	{xStart: 30, xEnd: 34, yStart: 17, yEnd: 27, obstacleType: Building},
	{xStart: 35, xEnd: 39, yStart: 17, yEnd: 24, obstacleType: Building},
	{xStart: 35, xEnd: 39, yStart: 25, yEnd: 27, obstacleType: Park},
	{xStart: 40, xEnd: 50, yStart: 24, yEnd: 27, obstacleType: Building},
	{xStart: 42, xEnd: 49, yStart: 28, yEnd: 31, obstacleType: Building},
	{xStart: 42, xEnd: 49, yStart: 32, yEnd: 42, obstacleType: Building},
	{xStart: 42, xEnd: 49, yStart: 43, yEnd: 50, obstacleType: Park},
	{xStart: 1, xEnd: 15, yStart: 31, yEnd: 40, obstacleType: Building},
	{xStart: 1, xEnd: 15, yStart: 41, yEnd: 44, obstacleType: Park},
	{xStart: 0, xEnd: 15, yStart: 45, yEnd: 50, obstacleType: Building},
	{xStart: 29, xEnd: 38, yStart: 31, yEnd: 41, obstacleType: Park},
	{xStart: 29, xEnd: 38, yStart: 42, yEnd: 44, obstacleType: Building},
	{xStart: 29, xEnd: 33, yStart: 45, yEnd: 49, obstacleType: Building},
	{xStart: 34, xEnd: 38, yStart: 45, yEnd: 49, obstacleType: Building},
	{xStart: 23, xEnd: 28, yStart: 31, yEnd: 41, obstacleType: Building},
	{xStart: 24, xEnd: 28, yStart: 42, yEnd: 46, obstacleType: Building},
	{xStart: 23, xEnd: 28, yStart: 47, yEnd: 50, obstacleType: Building},
	{xStart: 16, xEnd: 23, yStart: 42, yEnd: 50, obstacleType: Building},
	{xStart: 16, xEnd: 22, yStart: 31, yEnd: 34, obstacleType: Building},
	{xStart: 16, xEnd: 22, yStart: 35, yEnd: 41, obstacleType: Park},
}
