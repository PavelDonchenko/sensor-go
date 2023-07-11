package postgres

import (
	"embed"
	"fmt"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(fs embed.FS, cfg *config.Config) error {
	source, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("pgx://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)
	m, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		return err
	}
	err = m.Up()
	switch err {
	case nil:
		fmt.Println("Migration: uploaded successfully")
		return nil
	case migrate.ErrNoChange:
		fmt.Println("Migration: nothing to change")
		return nil
	default:
		return err
	}
}
