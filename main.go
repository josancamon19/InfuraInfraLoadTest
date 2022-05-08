package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"infura-test/endpoints"
	"log"
	"os"
)

func main() {

	app := fiber.New(fiber.Config{})
	app.Use(logger.New(logger.Config{Format: "[${time}] ${status} - ${latency} ${method} ${path}\n"}))
	app.Use(recover2.New())
	app.Get("/monitor", monitor.New())

	v1 := app.Group("/v1")
	v1.Get("/blockNumber", endpoints.GetBlockNumber)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("Defaulting to port: %s\n", port)
	}

	log.Fatal(app.Listen(":" + port))
}
