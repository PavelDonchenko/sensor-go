package dump_worker

import (
	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker struct {
	DB  *pgxpool.Pool
	log logging.Logger
	cfg config.Config
}
