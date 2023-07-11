package utils

import (
	"fmt"
	"os"

	"github.com/PavelDonchenko/sensor-go/config"
)

// ConnectionURLBuilder func for building URL connection.
func ConnectionURLBuilder(n string, cfg config.Config) (string, error) {
	// Define URL to connection.
	var url string

	// Switch given names.
	switch n {
	case "postgres":
		// URL for PostgreSQL connection.
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Username,
			cfg.Postgres.Password,
			cfg.Postgres.Database,
			"disable",
		)
	case "redis":
		// URL for Redis connection.
		url = fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		)
	case "fiber":
		// URL for Fiber connection.
		url = fmt.Sprintf(
			"%s:%s",
			cfg.HTTP.Host,
			cfg.HTTP.Port,
		)
	default:
		// Return error message.
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	// Return connection URL.
	return url, nil
}
