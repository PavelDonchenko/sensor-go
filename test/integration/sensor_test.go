package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/PavelDonchenko/sensor-go/internal/routes"
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
	err := SeedData(r.sensorStorage)
	assert.NoError(r.T(), err)

	defer func() {
		err := Truncate(r.sensorStorage)
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
			// Define a new Fiber app.
			app := fiber.New()

			// Define routes.
			routes.SensorRoute(app, r.sensorService)

			//w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/group/alpha/transparency/average", http.NoBody)

			resp, _ := app.Test(req, -1)
			fmt.Println(resp.Body)

			assert.Equal(r.T(), test.expectedStatusCode, resp.StatusCode)

		})
	}
}
