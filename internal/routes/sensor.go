package routes

import (
	"github.com/PavelDonchenko/sensor-go/internal/controllers"
	"github.com/PavelDonchenko/sensor-go/internal/service"
	"github.com/gofiber/fiber/v2"
)

func SensorRoute(a *fiber.App, service service.SensorService) {
	route := a.Group("/api/v1")

	route.Get("/group/:groupName/transparency/average", controllers.GetTransparency(service))
	route.Get("/group/:groupName/temperature/average", controllers.GetTemperature(service))
	route.Get("/group/:groupName/species", controllers.GetCurrentSpecies(service))
	route.Get("/group/:groupName/species/:top", controllers.GetCurrentTopSpecies(service))
}
