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
	app.Use(recover2.New())            // Used for handling Go panic()
	app.Get("/monitor", monitor.New()) // Visual tool for api performance

	v1 := app.Group("/v1")
	v1.Get("/:method", endpoints.InfuraHttpRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(app.Listen(":" + port))
}
