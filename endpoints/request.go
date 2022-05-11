package endpoints

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"infura-test/utils"
	"net/http"
	"strings"
	"time"
)

var availableMethods = []string{
	"eth_blockNumber", "eth_getTransactionByBlockNumberAndIndex",
}

func InfuraHttpRequest(c *fiber.Ctx) error {
	// Received path param baseUrl/:method
	method := c.Params("method")
	if method == "" {
		c.Status(http.StatusUnprocessableEntity)
		return c.JSON(fmt.Errorf("method in path required"))
	}

	// To make simpler the caching layer with redis, reducing the available methods to 2 works for the MVP
	if !contains(availableMethods, method) {
		c.Status(http.StatusUnprocessableEntity)
		return c.JSON(fmt.Errorf("we are not supporting method \"%s\" now", method))
	}

	// baseUrl/:method?params="param1,param2,param3" ~ extract the str value and cast it to a list
	strParams := c.Query("params")
	var params []string
	if strParams != "" {
		params = strings.Split(strParams, ",")
	}

	fmt.Printf("Calling method: \"%s\" with params: %v\n", method, params)
	redisKey := utils.GetRedisKeyFromInputs(method, params)

	// Try to find the cached value, as instead of making a request
	data, err := utils.RedisGetValue(redisKey)
	if err == nil && data["id"] != "" {
		fmt.Printf("Returning cached value for %s\n", redisKey)
		return c.JSON(data)
	}

	// Build the body for the Infura API request
	body, _ := utils.GetAPIBodyRequest(method, params)
	responseData, err := utils.InfuraAPIRequest(body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(err)
	}

	// Infura responses contain an "error" key with some specific details about the errors
	if responseData["error"] != nil {
		errorData := responseData["error"].(map[string]interface{})
		message := errorData["message"]
		c.Status(http.StatusInternalServerError)
		return c.JSON(message)
	}

	// For the method "eth_blockNumber" ~ new blocks are emitted every 12-14 seconds
	ttl := time.Second * 10
	if method == "eth_getTransactionByBlockNumberAndIndex" {
		// For "eth_getTransactionByBlockNumberAndIndex" ~ is getting a tx in an already appended block (immutable)
		// thus it'll never change. 0 value in redis means key has no ttl
		ttl = time.Second * 0
	}
	_ = utils.RedisSetKey(redisKey, responseData, ttl)
	return c.JSON(responseData)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
