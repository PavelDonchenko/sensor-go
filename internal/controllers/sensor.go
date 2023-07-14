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
