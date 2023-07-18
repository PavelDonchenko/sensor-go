package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/db"
	"github.com/PavelDonchenko/sensor-go/internal/handler"
	"github.com/PavelDonchenko/sensor-go/internal/service"
	"github.com/PavelDonchenko/sensor-go/internal/storage"
	"github.com/PavelDonchenko/sensor-go/pkg/cache"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
	"github.com/PavelDonchenko/sensor-go/pkg/postgres"
	"github.com/PavelDonchenko/sensor-go/pkg/utils"
	"github.com/PavelDonchenko/sensor-go/workers"
	"github.com/gofiber/fiber/v2"

	_ "github.com/PavelDonchenko/sensor-go/docs" // load API Docs files (Swagger)
)

// @title SENSOR API
// @version 1.0
// @description TEST API.
// @contact.email przmld033@gmail.com
// @BasePath /api
func main() {
	cfg := config.GetConfig("config.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.GetLogger()

	logger.Info("postgres initializing...")
	pool, err := postgres.NewClient(ctx, cfg)
	if err != nil {
		logger.Panic("error open postgres connection", err)
	}

	logger.Info("redis initializing...")
	redis, err := cache.NewCacheConn(*cfg)
	if err != nil {
		logger.Panic(err)
	}

	err = postgres.Migrate(db.Migrations, cfg)
	if err != nil {
		logger.Panic(err)
	}

	sensorStorage := storage.NewDatabase(pool, *cfg, logger)

	// in first running you must generate sensors and sensors group in PostgreSQL. checking if sensors are existing - everything ok,
	// if not - create them
	sensors, err := sensorStorage.GetAllSensors(ctx)
	if err != nil {
		logger.Panic(err)
	}

	if len(sensors) == 0 {
		logger.Info("Starting create new sensor and sensors group...")
		err = utils.GenerateSensors(ctx, *sensorStorage)
		if err != nil {
			logger.Panic("maybe need change config.yaml sensor_generated to true", err)
		}
	}

	worker := workers.NewWorker(ctx, sensorStorage, logger, *cfg)

	// worker is using to update sensor data
	logger.Info("Starting generate data for sensors...")
	go worker.Process()

	// Define a new Fiber app with config.
	app := fiber.New(fiber.Config{
		ReadTimeout: cfg.HTTP.ReadTimeOut,
	})

	sensorService := service.NewService(ctx, sensorStorage, logger, *cfg, redis)

	routes := handler.NewHandler(ctx, *cfg, sensorService)

	routes.Register(app)
	routes.RegisterSwagger(app)

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

	if err := a.Listen(fmt.Sprintf(":%s", cfg.HTTP.Port)); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
