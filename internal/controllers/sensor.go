package controllers

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/PavelDonchenko/sensor-go/internal/domain"
	"github.com/PavelDonchenko/sensor-go/internal/service"
	"github.com/PavelDonchenko/sensor-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// GetTransparency retrieves the transparency percentage for a sensor group.
//
// @Summary Get transparency percentage for a sensor group
// @Description Retrieves the transparency percentage for a sensor group based on the provided group name.
// @Tags sensors
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the sensor group"
// @Success 200 {number} number "transparency"
// @Failure 500 {string} string
// @Router /api/v1/group/{groupName}/transparency/average [get]
func GetTransparency(service service.SensorService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupName := c.Params("groupName")

		transparency, err := service.GetTransparency(c.Context(), strings.ToLower(groupName))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"error":          false,
			"msg":            nil,
			"transparency %": *transparency,
		})
	}
}

// GetTemperature retrieves the temperature in Celsius for a sensor group.
//
// @Summary Get temperature in Celsius for a sensor group
// @Description Retrieves the temperature in Celsius for a sensor group based on the provided group name.
// @Tags sensors
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the sensor group"
// @Success 200 {number} number "temperature"
// @Failure 500 {string} string
// @Router /api/v1/group/{groupName}/temperature/average [get]
func GetTemperature(service service.SensorService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupName := c.Params("groupName")

		temperature, err := service.GetTemperature(c.Context(), strings.ToLower(groupName))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"error":         false,
			"msg":           nil,
			"temperature C": math.Round(*temperature*1000) / 1000,
		})
	}
}

// GetCurrentSpecies retrieves the current detected fish species for a sensor group.
//
// @Summary Get current detected fish species for a sensor group
// @Description Retrieves the current detected fish species for a sensor group based on the provided group name.
// @Tags sensors
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the sensor group"
// @Success 200 {array} domain.ResponseDetectedFish
// @Failure 500 {string} string
// @Router /api/v1/group/{groupName}/species [get]
func GetCurrentSpecies(service service.SensorService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupName := c.Params("groupName")

		species, err := service.GetCurrentSpecies(c.Context(), strings.ToLower(groupName))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		var res []domain.ResponseDetectedFish

		for _, fish := range species {
			resFish := domain.ResponseDetectedFish{
				Name:  fish.Name,
				Count: fish.Count,
			}

			res = append(res, resFish)
		}

		return c.JSON(fiber.Map{
			"error":   false,
			"msg":     nil,
			"species": res,
		})
	}
}

// GetCurrentTopSpecies retrieves the current top detected fish species for a sensor group.
//
// @Summary Get current top detected fish species for a sensor group
// @Description Retrieves the current top detected fish species for a sensor group based on the provided group name and other optional parameters.
// @Tags sensors
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the sensor group"
// @Param top query integer true "Number of top species to retrieve"
// @Param from query string false "Start date for the period (UNIX timestamp)"
// @Param till query string false "End date for the period (UNIX timestamp)"
// @Success 200 {array} domain.ResponseDetectedFish
// @Failure 422 {string} string
// @Failure 500 {string} string
// @Router /api/v1/group/{groupName}/species/top/{top} [get]
func GetCurrentTopSpecies(service service.SensorService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupName := c.Params("groupName")

		stringTop := c.Params("top")

		top, err := strconv.Atoi(stringTop)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		start := c.Query("from")
		if start != "" {
			start = utils.ParseUnixToString(start)
		}

		end := c.Query("till")
		if end != "" {
			end = utils.ParseUnixToString(end)
		}

		species, err := service.GetCurrentTopSpecies(c.Context(), strings.ToLower(groupName), start, end, top)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		var res []domain.ResponseDetectedFish

		for _, fish := range species {
			resFish := domain.ResponseDetectedFish{
				Name:  fish.Name,
				Count: fish.Count,
			}

			res = append(res, resFish)
		}

		var msg string

		if start != "" {
			msg = fmt.Sprintf("species for period from %s till %s", start, end)
		}

		return c.JSON(fiber.Map{
			"error":   false,
			"msg":     msg,
			"species": res,
		})
	}
}
