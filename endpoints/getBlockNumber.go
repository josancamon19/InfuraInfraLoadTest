package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"infura-test/utils"
	"net/http"
)

func GetBlockNumber(c *fiber.Ctx) error {
	body, _ := utils.GetAPIBodyRequest("eth_blockNumber", []string{})
	responseData, err := utils.InfuraAPIRequest(body)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.JSON(err)
	}
	// TODO: handle Infura API Error in responseData
	return c.JSON(responseData)
}
