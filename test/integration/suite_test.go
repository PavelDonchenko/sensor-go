package test

import (
	"context"
	"log"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/db"
	"github.com/PavelDonchenko/sensor-go/internal/service"
	"github.com/PavelDonchenko/sensor-go/internal/storage"
	"github.com/PavelDonchenko/sensor-go/pkg/cache"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
	"github.com/PavelDonchenko/sensor-go/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	sensorService service.SensorService
	sensorStorage storage.Database
}

func (s *TestSuite) SetupTest() {
	logger := logging.GetLogger()
	cfg := &config.Config{}
	cfg.Postgres.Database = "test_sensor"
	cfg.Postgres.Port = "5432"
	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Password = "secret"
	cfg.Postgres.Username = "root"
	ctx := context.Background()

	pClient, err := pgxpool.New(ctx, "postgres://root:secret@localhost:5432/test_sensor?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db...")
	}

	redis, err := cache.NewCacheConn(*cfg)
	if err != nil {
		logger.Panic(err)
	}

	s.sensorStorage = *storage.NewDatabase(pClient, *cfg, logger)
	s.sensorService = service.NewService(ctx, &s.sensorStorage, logger, *cfg, redis)

	err = postgres.Migrate(db.Migrations, cfg)
	if err != nil {
		panic(err)
	}
}
