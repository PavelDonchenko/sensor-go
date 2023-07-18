package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SensorTestSuite struct {
	TestSuite
}

func TestSensorSuite(t *testing.T) {
	suite.Run(t, new(SensorTestSuite))
}

func (r *SensorTestSuite) TestGetTransparency() {
	err := SeedData(*r.sensorStorage)
	assert.NoError(r.T(), err)

	defer func() {
		err := Truncate(*r.sensorStorage)
		assert.NoError(r.T(), err)

	}()

	testCases := []struct {
		name               string
		groupName          string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			groupName:          "alpha",
			expectedStatusCode: 200,
		},
		{
			name:               "error wrong group name",
			groupName:          "wrong name",
			expectedStatusCode: 500},
	}

	for _, test := range testCases {
		r.Run(test.name, func() {
			app := fiber.New()

			url := fmt.Sprintf("/api/v1/group/%s/transparency/average", test.groupName)

			req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)

			r.handler.Register(app)

			resp, _ := app.Test(req, -1)

			assert.Equal(r.T(), test.expectedStatusCode, resp.StatusCode)
		})
	}
}

func (r *SensorTestSuite) TestGetTemperature() {
	err := SeedData(*r.sensorStorage)
	assert.NoError(r.T(), err)

	defer func() {
		err := Truncate(*r.sensorStorage)
		assert.NoError(r.T(), err)

	}()

	testCases := []struct {
		name               string
		groupName          string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			groupName:          "alpha",
			expectedStatusCode: 200,
		},
		{
			name:               "error wrong group name",
			groupName:          "wrong name",
			expectedStatusCode: 500},
	}

	for _, test := range testCases {
		r.Run(test.name, func() {
			app := fiber.New()

			url := fmt.Sprintf("/api/v1/group/%s/temperature/average", test.groupName)

			req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)

			r.handler.Register(app)

			resp, _ := app.Test(req, -1)

			assert.Equal(r.T(), test.expectedStatusCode, resp.StatusCode)
		})
	}
}

func (r *SensorTestSuite) TestGetCurrentSpecies() {
	err := SeedData(*r.sensorStorage)
	assert.NoError(r.T(), err)

	defer func() {
		err := Truncate(*r.sensorStorage)
		assert.NoError(r.T(), err)

	}()

	testCases := []struct {
		name               string
		groupName          string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			groupName:          "alpha",
			expectedStatusCode: 200,
		},
		{
			name:               "error wrong group name",
			groupName:          "wrong name",
			expectedStatusCode: 500},
	}

	for _, test := range testCases {
		r.Run(test.name, func() {
			app := fiber.New()

			url := fmt.Sprintf("/api/v1/group/%s/species", test.groupName)

			req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)

			r.handler.Register(app)

			resp, _ := app.Test(req, -1)

			assert.Equal(r.T(), test.expectedStatusCode, resp.StatusCode)
		})
	}
}

func (r *SensorTestSuite) TestGetTopSpecies() {
	err := SeedData(*r.sensorStorage)
	assert.NoError(r.T(), err)

	defer func() {
		err := Truncate(*r.sensorStorage)
		assert.NoError(r.T(), err)

	}()

	testCases := []struct {
		name               string
		groupName          string
		top                string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			groupName:          "alpha",
			top:                "3",
			expectedStatusCode: 200,
		},
		{
			name:               "error wrong top",
			groupName:          "alpha",
			top:                "wrong top",
			expectedStatusCode: 422,
		},
	}

	for _, test := range testCases {
		r.Run(test.name, func() {
			app := fiber.New()

			url := fmt.Sprintf("/api/v1/group/%s/species/top/%s", test.groupName, test.top)

			req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)

			r.handler.Register(app)

			resp, _ := app.Test(req, -1)

			assert.Equal(r.T(), test.expectedStatusCode, resp.StatusCode)
		})
	}
}
