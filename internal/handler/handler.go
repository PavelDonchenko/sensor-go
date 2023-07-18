package handler

import (
	"context"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/internal/service"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	ctx     context.Context
	cfg     config.Config
	service service.SensorService
}

func NewHandler(ctx context.Context, cfg config.Config, service service.SensorService) *Handler {
	return &Handler{ctx: ctx, cfg: cfg, service: service}
}

func (h *Handler) Register(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Get("/group/:groupName/transparency/average", h.GetTransparency)
	route.Get("/group/:groupName/temperature/average", h.GetTemperature)
	route.Get("/group/:groupName/species", h.GetCurrentSpecies)
	route.Get("/group/:groupName/species/top/:top", h.GetCurrentTopSpecies)
	route.Get("/region/temperature/min", h.GeRegionMinTemperature)
	route.Get("/region/temperature/max", h.GeRegionMaxTemperature)
	route.Get("/sensor/:codename/temperature/average", h.GetAverageSensorTemperature)
}
func (h *Handler) RegisterSwagger(a *fiber.App) {
	// Create routes group.
	route := a.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.HandlerDefault)
}
