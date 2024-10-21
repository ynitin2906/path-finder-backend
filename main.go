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

// Directions for left-right and up-down movements
var horizontalDirections = []Point{
	{0, 1},  // right
	{0, -1}, // left
}

var verticalDirections = []Point{
	{1, 0},  // down
	{-1, 0}, // up
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://path-finder-frontend-delta.vercel.app",
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

	var dfsHelper func(Point, bool) bool
	dfsHelper = func(curr Point, lastWasVertical bool) bool {
		// Check boundaries and if already visited
		if curr.X < 0 || curr.X >= 20 || curr.Y < 0 || curr.Y >= 20 || visited[curr] {
			return false
		}
		// Check if we reached the end
		if curr == end {
			path = append(path, []int{curr.X, curr.Y})
			return true
		}

		// Mark the node as visited
		visited[curr] = true
		path = append(path, []int{curr.X, curr.Y})

		// Explore horizontal movements
		for _, dir := range horizontalDirections {
			next := Point{curr.X + dir.X, curr.Y + dir.Y}
			if dfsHelper(next, false) { // After horizontal move, set lastWasVertical to false
				return true // Exit if path found
			}
		}

		// Explore vertical movements only if the last move was horizontal
		if !lastWasVertical {
			for _, dir := range verticalDirections {
				next := Point{curr.X + dir.X, curr.Y + dir.Y}
				if dfsHelper(next, true) { // After vertical move, set lastWasVertical to true
					return true // Exit if path found
				}
			}
		}

		// Backtrack if no solution found
		visited[curr] = false // Unmark the node
		path = path[:len(path)-1]
		return false
	}

	dfsHelper(start, false)

	// Check if the last point in path is the endpoint
	if len(path) > 0 && path[len(path)-1][0] == end.X && path[len(path)-1][1] == end.Y {
		return path // Return only if the last point is the endpoint
	}
	return nil // Return nil if no path was found
}
