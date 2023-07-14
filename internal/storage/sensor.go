package storage

import (
	"context"
	"math/rand"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/PavelDonchenko/sensor-go/internal/domain"
	"github.com/PavelDonchenko/sensor-go/pkg/generations"
	"github.com/PavelDonchenko/sensor-go/pkg/logging"
	"github.com/PavelDonchenko/sensor-go/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorPostgres interface {
	SaveDetectedFish(ctx context.Context, fish domain.DetectedFish) (*domain.DetectedFish, error)
	UpdateSensorData(ctx context.Context, sensor domain.Sensor) error
	GetAllSensors(ctx context.Context) ([]domain.Sensor, error)
	GetTransparency(ctx context.Context, groupName string) (float64, error)
	GetTemperature(ctx context.Context, groupName string) (float64, error)
	GetSpecies(ctx context.Context, groupName string) ([]domain.DetectedFish, error)
	GetTopSpecies(ctx context.Context, groupName, start, end string, top int) ([]domain.DetectedFish, error)
}

type Database struct {
	DB  *pgxpool.Pool
	Cfg config.Config
	log logging.Logger
}

func NewDatabase(DB *pgxpool.Pool, cfg config.Config, log logging.Logger) *Database {
	return &Database{DB: DB, Cfg: cfg, log: log}
}

func (d *Database) CreateSensorGroup(ctx context.Context, name string, id int) error {
	query := "INSERT INTO sensor_group (id, name) VALUES ($1, $2)"
	_, err := d.DB.Exec(ctx, query, id, name)
	if err != nil {
		err = postgres.ErrDoQuery(err)
		d.log.Error(err)
		return err
	}

	return nil
}

func (d *Database) CreateSensorsForGroup(ctx context.Context, group []string, sensorCount int) error {
	tx, err := d.DB.Begin(ctx)
	if err != nil {
		err = postgres.ErrCreateTx(err)
		d.log.Error(err)
		return err
	}

	for groupID, groupName := range group {
		for i := 1; i <= sensorCount; i++ {
			// Generate random coordinates within the group's range
			coordinates := generations.GenerateCoordinates(groupID)

			codename := domain.Codename{
				Name:          groupName,
				SensorGroupID: i,
			}

			sensor := &domain.Sensor{
				Codename:       codename,
				Coordinates:    coordinates,
				Transparency:   rand.Intn(101),
				Temperature:    generations.GenerateTemperature(coordinates.Z),
				DataOutputRate: generations.GenerateRandomInt(), // Random data output rate between 5, 10, 15, 20, 25
			}

			query := "INSERT INTO sensor (group_id, group_name, in_group_id, data_output_rate, x, y, z, transparency, temperature) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

			_, err = tx.Exec(ctx, query, groupID, sensor.Codename.Name, sensor.Codename.SensorGroupID, sensor.DataOutputRate, sensor.Coordinates.X, sensor.Coordinates.Y, sensor.Coordinates.Z, sensor.Transparency, sensor.Temperature)
			if err != nil {
				err = postgres.ErrExecQuery(err)
				d.log.Error(err)
				_ = tx.Rollback(ctx)
				return err
			}

		}
	}
	updateQuery := `UPDATE sensor_group 
							SET sensors = (
							SELECT ARRAY_AGG(id) 
							FROM sensor 
							WHERE sensor.group_name = sensor_group.name)
                    		WHERE EXISTS (SELECT 1 FROM sensor WHERE sensor.group_name = sensor_group.name)`

	_, err = tx.Exec(ctx, updateQuery)
	if err != nil {
		err = postgres.ErrExecQuery(err)
		d.log.Error(err)
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = postgres.ErrCommit(err)
		d.log.Error(err)
		return err
	}

	return nil
}

func (d *Database) SaveDetectedFish(ctx context.Context, fish domain.DetectedFish) (*domain.DetectedFish, error) {
	var detectedFish domain.DetectedFish
	query := "INSERT INTO detected_fish (name, count, sensorid) VALUES ($1, $2, $3) RETURNING id, name, count"
	err := d.DB.QueryRow(ctx, query, fish.Name, fish.Count, fish.SensorID).Scan(&detectedFish.ID, &detectedFish.Name, &detectedFish.Count)
	if err != nil {
		err = postgres.ErrScan(err)
		d.log.Error(err)
		return nil, err
	}
	return &detectedFish, nil
}

func (d *Database) UpdateSensorData(ctx context.Context, sensor domain.Sensor) error {
	var fishesID []uuid.UUID
	for _, fish := range sensor.DetectedFish {
		fishesID = append(fishesID, fish.ID)
	}

	sensorQuery := "UPDATE sensor SET transparency = $1, temperature = $2, fishes = $3 WHERE id = $4"

	_, err := d.DB.Exec(ctx, sensorQuery, sensor.Transparency, sensor.Temperature, fishesID, sensor.ID)
	if err != nil {
		err = postgres.ErrExecQuery(err)
		d.log.Error(err)
		return err
	}

	//for _, fish := range sensor.DetectedFish {
	//	_, err = d.SaveDetectedFish(ctx, fish)
	//	if err != nil {
	//		return fmt.Errorf("error save fetected fish for sensor: %s, err: %v", sensor.ID, err)
	//	}
	//}

	return nil
}

func (d *Database) GetAllSensors(ctx context.Context) ([]domain.Sensor, error) {
	query := "SELECT id, temperature, transparency, created_at, in_group_id, group_name, data_output_rate, x, y, z FROM sensor"

	rows, err := d.DB.Query(ctx, query)
	if err != nil {
		err = postgres.ErrDoQuery(err)
		d.log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var sensors []domain.Sensor

	for rows.Next() {
		var sensor domain.Sensor
		err := rows.Scan(
			&sensor.ID,
			&sensor.Temperature,
			&sensor.Transparency,
			&sensor.CreatedAt,
			&sensor.Codename.SensorGroupID,
			&sensor.Codename.Name,
			&sensor.DataOutputRate,
			&sensor.Coordinates.X,
			&sensor.Coordinates.Y,
			&sensor.Coordinates.Z,
		)
		if err != nil {
			err = postgres.ErrScan(err)
			d.log.Error(err)
			return nil, err
		}

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

func (d *Database) GetTransparency(ctx context.Context, groupName string) (float64, error) {
	query := "SELECT AVG(transparency) from sensor WHERE group_name = $1"

	var average float64

	err := d.DB.QueryRow(ctx, query, groupName).Scan(&average)
	if err != nil {
		err = postgres.ErrScan(err)
		d.log.Error(err)
		return 0, err
	}

	return average, nil
}

func (d *Database) GetTemperature(ctx context.Context, groupName string) (float64, error) {
	query := "SELECT AVG(temperature) from sensor WHERE group_name = $1"

	var average float64

	err := d.DB.QueryRow(ctx, query, groupName).Scan(&average)
	if err != nil {
		err = postgres.ErrScan(err)
		d.log.Error(err)
		return 0, err
	}

	return average, nil
}

func (d *Database) GetSpecies(ctx context.Context, groupName string) ([]domain.DetectedFish, error) {
	query := `SELECT df.name, SUM(df.count) AS total_count 
		 FROM detected_fish df 
         WHERE df.id IN (
		 SELECT DISTINCT df.id 
		 FROM detected_fish df 
		 JOIN sensor s ON df.id = ANY(s.fishes) AND s.group_name = $1)
		 GROUP BY df.name`

	var fishes []domain.DetectedFish

	rows, err := d.DB.Query(ctx, query, groupName)
	if err != nil {
		err = postgres.ErrDoQuery(err)
		d.log.Error(err)
		return nil, err
	}

	for rows.Next() {
		var fish domain.DetectedFish
		err := rows.Scan(
			&fish.Name,
			&fish.Count,
		)
		if err != nil {
			err = postgres.ErrScan(err)
			d.log.Error(err)
			return nil, err
		}

		fishes = append(fishes, fish)
	}

	return fishes, nil
}

func (d *Database) GetTopSpecies(ctx context.Context, groupName, start, end string, top int) ([]domain.DetectedFish, error) {
	var query string

	if start == "" {
		query = `SELECT df.name, SUM(df.count) AS total_count 
		 FROM detected_fish df 
         WHERE df.id IN (
		 SELECT DISTINCT df.id 
		 FROM detected_fish df 
		 JOIN sensor s ON df.id = ANY(s.fishes) AND s.group_name = $1)
		 GROUP BY df.name
		 ORDER BY total_count DESC
		 LIMIT $2`
	} else {
		query = `SELECT df.name, SUM(df.count) AS total_count  FROM detected_fish as df
         JOIN sensor s on s.id = df.sensorid
         WHERE group_name = $1 AND s.created_at between $3 AND $4
         GROUP BY df.name
         ORDER BY total_count DESC
         LIMIT $2`
	}

	var fishes []domain.DetectedFish

	var err error
	var rows pgx.Rows

	if start == "" {
		rows, err = d.DB.Query(ctx, query, groupName, top)
	} else {
		rows, err = d.DB.Query(ctx, query, groupName, top, start, end)
	}
	if err != nil {
		err = postgres.ErrDoQuery(err)
		d.log.Error(err)
		return nil, err
	}

	for rows.Next() {
		var fish domain.DetectedFish
		err := rows.Scan(
			&fish.Name,
			&fish.Count,
		)
		if err != nil {
			err = postgres.ErrScan(err)
			d.log.Error(err)
			return nil, err
		}

		fishes = append(fishes, fish)
	}

	return fishes, nil
}
