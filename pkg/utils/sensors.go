package utils

import (
	"context"
	"strings"

	"github.com/PavelDonchenko/sensor-go/internal/storage"
)

func GenerateSensors(ctx context.Context, db storage.Database) error {
	groupNames := strings.Split(db.Cfg.GroupNames, " ")

	for i, name := range groupNames {
		err := db.CreateSensorGroup(ctx, name, i)
		if err != nil {
			return err
		}
	}

	err := db.CreateSensorsForGroup(ctx, groupNames, db.Cfg.CountSensorInGroup)
	if err != nil {
		return err
	}
	return nil
}
