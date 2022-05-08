package endpoints

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"infura-test/utils"
	"net/http"
	"strings"
)

func InfuraHttpRequest(c *fiber.Ctx) error {
	method := c.Params("method")
	if method == "" {
		c.Status(http.StatusUnprocessableEntity)
		return c.JSON(fmt.Errorf("method in path required"))
	}

	strParams := c.Query("params")
	var params []string
	if strParams != "" {
		params = strings.Split(strParams, ",")
	}

	fmt.Printf("Calling method: \"%s\" with params: %v\n", method, params)

	body, _ := utils.GetAPIBodyRequest(method, params)
	responseData, err := utils.InfuraAPIRequest(body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(err)
	}

	if responseData["error"] != nil {
		errorData := responseData["error"].(map[string]interface{})
		message := errorData["message"]
		fmt.Printf("Infura API Error response: %v\n", message)

		c.Status(http.StatusInternalServerError)
		return c.JSON(message)
	}

	return c.JSON(responseData)
}
