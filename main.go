package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Request struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

var directions = []Point{
	{0, 1}, {1, 0}, {0, -1}, {-1, 0},
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Post("/find-path", func(c *fiber.Ctx) error {
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		path := dfs(req.Start, req.End)
		if path == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Path not found"})
		}
		return c.JSON(fiber.Map{"path": path})
	})

	log.Fatal(app.Listen(":8080"))
}

func dfs(start, end Point) [][]int {
	visited := make(map[Point]bool)
	var path [][]int
	found := false

	var dfsHelper func(Point) bool
	dfsHelper = func(curr Point) bool {
		if curr.X < 0 || curr.X >= 20 || curr.Y < 0 || curr.Y >= 20 || visited[curr] {
			return false
		}
		visited[curr] = true
		path = append(path, []int{curr.X, curr.Y})

		if curr == end {
			found = true
			return true
		}

		for _, dir := range directions {
			next := Point{curr.X + dir.X, curr.Y + dir.Y}
			if dfsHelper(next) {
				return true
			}
		}

		// Backtrack if no solution found
		if !found {
			path = path[:len(path)-1]
		}
		return false
	}

	dfsHelper(start)
	if found {
		return path
	}
	return nil // Return nil if no path was found
}
