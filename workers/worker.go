package workers

import (
	"context"
	"math/rand"
	"time"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/internal/domain"
	"github.com/PavelDonchenko/sensor-go/internal/storage"
	"github.com/PavelDonchenko/sensor-go/pkg/generations"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
)

type Worker struct {
	DB        storage.SensorPostgres
	ctx       context.Context
	log       logging.Logger
	cfg       config.Config
	errorChan chan error
}

func NewWorker(ctx context.Context, DB storage.SensorPostgres, log logging.Logger, cfg config.Config) *Worker {
	errChan := make(chan error)
	return &Worker{
		DB:        DB,
		ctx:       ctx,
		log:       log,
		cfg:       cfg,
		errorChan: errChan,
	}
}

func (w *Worker) Process() {
	sensors, err := w.DB.GetAllSensors(w.ctx)
	if err != nil {
		w.log.Error(err)
		w.errorChan <- err
	}

	for _, sensor := range sensors {
		go w.generateSensorData(sensor)
	}
}

func (w *Worker) generateSensorData(sensor domain.Sensor) {
	duration := time.Duration(sensor.DataOutputRate) * time.Second

	ticker := time.NewTicker(duration)

	for {
		select {
		case <-ticker.C:
			fishes, err := w.generateFishData(sensor)
			if err != nil {
				w.log.Error(err)
				w.errorChan <- err
			}

			temperature := generations.GenerateTemperature(sensor.Coordinates.Z)

			err = w.DB.SaveTemperature(w.ctx, temperature, sensor.ID)
			if err != nil {
				w.errorChan <- err
			}

			transparency := generations.GenerateTransparency(sensor.Transparency)

			toUpdate := domain.Sensor{
				ID:           sensor.ID,
				Temperature:  temperature,
				Transparency: transparency,
				DetectedFish: fishes,
				UpdatedAt:    time.Now(),
			}

			err = w.DB.UpdateSensorData(w.ctx, toUpdate)
			if err != nil {
				w.log.Error(err)
				w.errorChan <- err
			}
		case <-w.errorChan:
			return
		}
	}
}

func (w *Worker) generateFishData(sensor domain.Sensor) ([]domain.DetectedFish, error) {
	fishSpecies := []string{"Atlantic Cod", "Sailfish", "Tuna", "Salmon", "Marlin", "Barracuda"}

	var detectedFish []domain.DetectedFish

	numFish := rand.Intn(len(fishSpecies))

	uniqueFish := make(map[string]bool)

	for len(uniqueFish) < numFish {
		fish := fishSpecies[rand.Intn(len(fishSpecies))]
		if !uniqueFish[fish] {
			count := rand.Intn(20) + 1
			detectedFish = append(detectedFish, domain.DetectedFish{Name: fish, Count: count})
			uniqueFish[fish] = true
		}
	}

	var detectedFishes []domain.DetectedFish

	for _, fish := range detectedFish {
		fish.SensorID = sensor.ID
		detectedFish, err := w.DB.SaveDetectedFish(w.ctx, fish)
		if err != nil {
			return nil, err
		}

		detectedFishes = append(detectedFishes, *detectedFish)
	}

	return detectedFishes, nil
}
