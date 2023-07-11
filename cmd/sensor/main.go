package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
	"github.com/PavelDonchenko/sensor-go/pkg/postgres"
	"github.com/PavelDonchenko/sensor-go/pkg/utils"
	"github.com/gofiber/fiber/v2"

	_ "github.com/PavelDonchenko/sensor-go/docs" // load API Docs files (Swagger)
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Define config.
	cfg := config.GetConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.GetLogger()

	_, err := postgres.NewClient(ctx, cfg)
	if err != nil {
		logger.Panic("error open postgres connection")
	}

	// Define a new Fiber app with config.
	app := fiber.New(fiber.Config{
		ReadTimeout: cfg.HTTP.ReadTimeOut,
	})

	//// Middlewares.
	//middleware.FiberMiddleware(app) // Register Fiber's middleware for app.
	//
	//// Routes.
	//routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	//routes.PublicRoutes(app)  // Register a public routes for app.
	//routes.PrivateRoutes(app) // Register a private routes for app.
	//routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server with graceful shutdown.
	StartServerWithGracefulShutdown(app, *cfg)
}

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(a *fiber.App, cfg config.Config) {
	// Create channel for idle connections.
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt) // Catch OS signals.
		<-sigint

		// Received an interrupt signal, shutdown.
		if err := a.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	// Build Fiber connection URL.
	fiberConnURL, _ := utils.ConnectionURLBuilder("fiber", cfg)

	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
