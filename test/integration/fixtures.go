package test

import (
	"context"
	"fmt"

	"github.com/PavelDonchenko/sensor-go/internal/storage"
	"github.com/PavelDonchenko/sensor-go/pkg/utils"
)

func SeedData(db storage.Database) error {
	ctx := context.Background()
	err := utils.GenerateSensors(ctx, db)
	if err != nil {
		return err
	}

	sensors, err := db.GetAllSensors(ctx)
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		err = db.UpdateSensorData(ctx, sensor)
		if err != nil {
			fmt.Print(err)
			return err
		}
	}

	return nil
}

func Truncate(db storage.Database) error {
	_, err := db.DB.Exec(context.Background(), "TRUNCATE table sensor_group, sensor, detected_fish")
	if err != nil {
		return err
	}

	return nil
}
