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
	app := fiber.New()
	app.Use(logger.New(logger.Config{Format: "[${time}] ${status} - ${latency} ${method} ${path}\n"}))
	app.Use(recover2.New()) // Used for handling Go panic()

	app.Get("/monitor", monitor.New()) // Visual tool for api performance
	app.Get("/:method", endpoints.InfuraHttpRequest)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
