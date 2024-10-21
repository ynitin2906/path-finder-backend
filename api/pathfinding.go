package api

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Node struct {
	Row      int  `json:"row"`
	Col      int  `json:"col"`
	IsWall   bool `json:"isWall"`
	IsStart  bool `json:"isStart"`
	IsFinish bool `json:"isFinish"`
}

func GetPathfinding(c *fiber.Ctx) error {
	rand.Seed(time.Now().UnixNano())
	grid := make([][]Node, 25)
	for i := range grid {
		grid[i] = make([]Node, 35)
		for j := range grid[i] {
			grid[i][j] = Node{
				Row:      i,
				Col:      j,
				IsWall:   rand.Intn(5) == 0, // 20% chance to create a wall
				IsStart:  i == 5 && j == 5,
				IsFinish: i == 5 && j == 15,
			}
		}
	}
	return c.JSON(grid)
}
