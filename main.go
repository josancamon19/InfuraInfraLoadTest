package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"infura-test/endpoints"
	"infura-test/utils"
	"log"
	"os"
	"time"
)

func initRedis() {
	var dbAddress = os.Getenv("REDIS_DB_ADDRESS")
	if dbAddress == "" {
		dbAddress = "localhost:6379"
	}
	utils.RDB = redis.NewClient(&redis.Options{Addr: dbAddress, Password: "", DB: 0})
}
func main() {

	initRedis()

	app := fiber.New(fiber.Config{})
	app.Use(favicon.New())
	app.Use(logger.New(logger.Config{Format: "[${time}] ${status} - ${latency} ${method} ${path}\n"}))
	app.Use(recover.New())             // Used for handling Go panic()
	app.Get("/monitor", monitor.New()) // Visual tool for api performance

	// How's this caching layer compared to the redis one?
	// Details in ./report/Report.mdx
	app.Use(cache.New(cache.Config{Expiration: 10 * time.Second}))

	app.Get("/:method", endpoints.InfuraHttpRequest)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
