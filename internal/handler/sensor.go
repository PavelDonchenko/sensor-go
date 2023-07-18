package handler

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/PavelDonchenko/sensor-go/internal/domain"
	"github.com/PavelDonchenko/sensor-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// GetTransparency retrieves the transparency percentage for a sensor group.
//
// @Summary Get transparency percentage for a sensor group
// @Description Retrieves the transparency percentage for a sensor group based on the provided group name.
// @Tags group
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the sensor group"
// @Success 200 {number} number "transparency"
// @Failure 500 {string} string
// @Router /api/v1/group/{groupName}/transparency/average [get]
func (h *Handler) GetTransparency(c *fiber.Ctx) error {
	groupName := c.Params("groupName")

	transparency, err := h.service.GetTransparency(h.ctx, strings.ToLower(groupName))
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

// GetTemperature retrieves the temperature in Celsius for a sensor group.
//
// @Summary Get temperature in Celsius for a sensor group
// @Description Retrieves the temperature in Celsius for a sensor group based on the provided group name.
// @Tags group
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the sensor group"
// @Success 200 {number} number "temperature"
// @Failure 500 {string} string
// @Router /api/v1/group/{groupName}/temperature/average [get]
func (h *Handler) GetTemperature(c *fiber.Ctx) error {
	groupName := c.Params("groupName")

	temperature, err := h.service.GetTemperature(h.ctx, strings.ToLower(groupName))
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

// GetCurrentSpecies retrieves the current detected fish species for a sensor group.
//
// @Summary Get current detected fish species for a sensor group
// @Description Retrieves the current detected fish species for a sensor group based on the provided group name.
// @Tags group
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the sensor group"
// @Success 200 {array} domain.ResponseDetectedFish
// @Failure 500 {string} string
// @Router /api/v1/group/{groupName}/species [get]
func (h *Handler) GetCurrentSpecies(c *fiber.Ctx) error {
	groupName := c.Params("groupName")

	species, err := h.service.GetCurrentSpecies(h.ctx, strings.ToLower(groupName))
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

// GetCurrentTopSpecies retrieves the current top detected fish species for a sensor group.
//
// @Summary Get current top detected fish species for a sensor group
// @Description Retrieves the current top detected fish species for a sensor group based on the provided group name and other optional parameters.
// @Tags group
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
func (h *Handler) GetCurrentTopSpecies(c *fiber.Ctx) error {
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

	species, err := h.service.GetCurrentTopSpecies(h.ctx, strings.ToLower(groupName), start, end, top)
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

// GeRegionMinTemperature retrieves the current minimum temperature according to region.
//
// @Summary Get current minimum temperature according to region.
// @Description Retrieves the current minimum temperature with optional parameters.
// @Tags region
// @Accept json
// @Produce json
// @Param xMin query number true  "minimum X coordinate"
// @Param xMax query number true	"maximum X coordinate"
// @Param yMax query number true	"maximum Y coordinate"
// @Param yMix query number true	"minimum Y coordinate"
// @Param zMin query number true	"minimum Z coordinate"
// @Param zMax query number true	"maximum Z coordinate"
// @Success 200 {number} number
// @Failure 500 {string} string
// @Router /api/v1/region/temperature/min [get]
func (h *Handler) GeRegionMinTemperature(c *fiber.Ctx) error {
	xMin, _ := strconv.ParseFloat(c.Query("xMin"), 64)
	xMax, _ := strconv.ParseFloat(c.Query("xMax"), 64)
	yMin, _ := strconv.ParseFloat(c.Query("yMin"), 64)
	yMax, _ := strconv.ParseFloat(c.Query("yMax"), 64)
	zMin, _ := strconv.ParseFloat(c.Query("zMin"), 64)
	zMax, _ := strconv.ParseFloat(c.Query("zMax"), 64)

	region := domain.Region{
		XMin: xMin,
		XMax: xMax,
		YMin: yMin,
		YMax: yMax,
		ZMin: zMin,
		ZMax: zMax,
	}

	flag := "MIN"

	temperature, err := h.service.GetRegionTemperature(h.ctx, region, flag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":              false,
		"msg":                nil,
		"region temperature": temperature,
	})
}

// GeRegionMaxTemperature retrieves the current maximum temperature according to region.
//
// @Summary Get current maximum temperature according to region.
// @Description Retrieves the current maximum temperature with optional parameters.
// @Tags region
// @Accept json
// @Produce json
// @Param xMin query number true  "minimum X coordinate"
// @Param xMax query number true	"maximum X coordinate"
// @Param yMax query number true	"maximum Y coordinate"
// @Param yMix query number true	"minimum Y coordinate"
// @Param zMin query number true	"minimum Z coordinate"
// @Param zMax query number true	"maximum Z coordinate"
// @Success 200 {number} number
// @Failure 500 {string} string
// @Router /api/v1/region/temperature/max [get]
func (h *Handler) GeRegionMaxTemperature(c *fiber.Ctx) error {
	xMin, _ := strconv.ParseFloat(c.Query("xMin"), 64)
	xMax, _ := strconv.ParseFloat(c.Query("xMax"), 64)
	yMin, _ := strconv.ParseFloat(c.Query("yMin"), 64)
	yMax, _ := strconv.ParseFloat(c.Query("yMax"), 64)
	zMin, _ := strconv.ParseFloat(c.Query("zMin"), 64)
	zMax, _ := strconv.ParseFloat(c.Query("zMax"), 64)

	region := domain.Region{
		XMin: xMin,
		XMax: xMax,
		YMin: yMin,
		YMax: yMax,
		ZMin: zMin,
		ZMax: zMax,
	}

	flag := "MAX"

	temperature, err := h.service.GetRegionTemperature(h.ctx, region, flag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":              false,
		"msg":                nil,
		"region temperature": temperature,
	})
}

// GetAverageSensorTemperature retrieves the current top detected fish species for a sensor group.
//
// @Summary Get average temperature from sensor
// @Description Retrieves the average temperature based on the  optional parameters.
// @Tags sensor
// @Accept json
// @Produce json
// @Param codename path string true "name of the group and id inside the group"
// @Param from query string false "Start date for the period (UNIX timestamp)"
// @Param till query string false "End date for the period (UNIX timestamp)"
// @Success 200 {number} number
// @Failure 500 {string} string
// @Router /api/v1/sensor/{codename}/temperature/average [get]
func (h *Handler) GetAverageSensorTemperature(c *fiber.Ctx) error {
	codename := c.Params("codename")

	group, inGroupID := utils.ParseCodename(codename)

	start := c.Query("from")
	if start != "" {
		start = utils.ParseUnixToString(start)
	}

	end := c.Query("till")
	if end != "" {
		end = utils.ParseUnixToString(end)
	}

	temperature, err := h.service.GetSensorTemperature(h.ctx, inGroupID, group, start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":       false,
		"msg":         fmt.Sprintf("temperature from scanner %s, period from %s till %s", codename, start, end),
		"temperature": math.Round(*temperature*1000) / 1000,
	})
}
