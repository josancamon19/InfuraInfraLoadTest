package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	app := fiber.New()
	//app.Use(logger.New(logger.Config{Format: "[${time}] ${status} - ${latency} ${method} ${path}\n"}))

	v1 := app.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {

		body := map[string]interface{}{
			"jsonrpc": "2.0", "method": "eth_blockNumber", "params": []string{}, "id": 1,
		}
		jsonData, err := json.Marshal(body)
		requestUrl := os.Getenv("API_BASE_URL") + os.Getenv("INFURA_PROJECT_ID")
		response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err)
		}

		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(string(responseBody))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("Defaulting to port: %s\n", port)
	}

	log.Fatal(app.Listen(":" + port))
}
