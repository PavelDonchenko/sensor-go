package test

import (
	"context"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/db"
	"github.com/PavelDonchenko/sensor-go/internal/handler"
	"github.com/PavelDonchenko/sensor-go/internal/service"
	"github.com/PavelDonchenko/sensor-go/internal/storage"
	"github.com/PavelDonchenko/sensor-go/pkg/cache"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
	"github.com/PavelDonchenko/sensor-go/pkg/postgres"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	sensorService service.SensorService
	sensorStorage *storage.Database
	handler       *handler.Handler
}

func (s *TestSuite) SetupTest() {
	logger := logging.GetLogger()

	cfg := config.GetConfig("../../config.yaml")

	cfg.Postgres.Database = "test_sensor"
	cfg.Postgres.Host = "localhost"
	cfg.Redis.Address = "localhost:6379"

	ctx := context.Background()

	pClient, err := postgres.NewClient(ctx, cfg)
	if err != nil {
		logger.Panic("error open postgres connection", err)
	}

	redis, err := cache.NewCacheConn(*cfg)
	if err != nil {
		logger.Panic(err)
	}

	s.sensorStorage = storage.NewDatabase(pClient, *cfg, logger)

	s.sensorService = service.NewService(ctx, s.sensorStorage, logger, *cfg, redis)

	s.handler = handler.NewHandler(ctx, *cfg, s.sensorService)

	err = postgres.Migrate(db.Migrations, cfg)
	if err != nil {
		logger.Panic("error migration", err)
	}
}
