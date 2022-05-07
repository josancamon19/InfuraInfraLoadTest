package main

import (
	"github.com/gofiber/fiber/v2"
	"infura-test/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	app := fiber.New()
	//app.Use(logger.New(logger.Config{Format: "[${time}] ${status} - ${latency} ${method} ${path}\n"}))

	v1 := app.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {

		body, err := utils.GetAPIBodyRequest("eth_accounts", []string{})
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(err)
		}

		responseData, err := utils.InfuraAPIRequest(body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(err)
		}

		return c.JSON(responseData)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("Defaulting to port: %s\n", port)
	}

	log.Fatal(app.Listen(":" + port))
}
