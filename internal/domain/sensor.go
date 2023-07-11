package domain

import "time"

type Sensor struct {
	ID              int          `json:"id"`
	SensorGroupID   int          `json:"sensor_group_id"`
	DataOutputRate  int          `json:"dataOutput_rate"`
	Temperature     float64      `json:"temperature"`
	Transparency    int          `json:"transparency"`
	LastMeasurement time.Time    `json:"last_measurement"`
	Codename        Codename     `json:"codename"`
	Coordinates     Coordinates  `json:"coordinates"`
	DetectedFish    DetectedFish `json:"detected_fish"`
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
	SensorID int    `json:"sensor_id"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

type Codename struct {
	Name          string `json:"name"`
	SensorGroupID int    `json:"sensor_group_id"`
}
