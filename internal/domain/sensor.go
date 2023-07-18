package domain

import (
	"time"

	"github.com/google/uuid"
)

type Sensor struct {
	ID             uuid.UUID      `json:"id"`
	DataOutputRate int            `json:"data_output_rate"`
	Temperature    float64        `json:"temperature"`
	Transparency   int            `json:"transparency"`
	Codename       Codename       `json:"codename"`
	Coordinates    Coordinates    `json:"coordinates"`
	DetectedFish   []DetectedFish `json:"detected_fish"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CreatedAt      time.Time      `json:"created_at"`
}

type SensorGroup struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Sensors []Sensor `json:"sensors"`
}

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type DetectedFish struct {
	ID       uuid.UUID `json:"id,omitempty"`
	SensorID uuid.UUID `json:"sensor_id,omitempty"`
	Name     string    `json:"name"`
	Count    int       `json:"count"`
}

type ResponseDetectedFish struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Codename struct {
	Name          string `json:"name"`
	SensorGroupID int    `json:"sensor_group_id"`
}

type Region struct {
	XMin float64 `yaml:"x_min"`
	XMax float64 `yaml:"x_max"`
	YMin float64 `yaml:"y_min"`
	YMax float64 `yaml:"y_max"`
	ZMin float64 `yaml:"z_min"`
	ZMax float64 `yaml:"z_max"`
}
