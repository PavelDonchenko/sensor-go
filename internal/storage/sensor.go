package storage

import (
	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	DB  *pgxpool.Pool
	cfg config.Config
}
