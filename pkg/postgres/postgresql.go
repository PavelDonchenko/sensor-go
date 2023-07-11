package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewClient Create Postgres pgx connection with attempts
func NewClient(ctx context.Context, cfg *config.Config) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)

	if len(dsn) < 40 {
		return nil, fmt.Errorf("wrong connection sring")
	}

	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			fmt.Println("failed to connect to postgesql... Going to do the next attempt")
			return err
		}
		return nil
	}, cfg.Postgres.MaxAttempts, 5*time.Second)
	if err != nil {
		log.Fatalln("All attempts are exceeded. Unable to connect to postgres")
	}
	return pool, nil
}

// DoWithTries  provide attempts to connect db
func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}
