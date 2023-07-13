package generations

import (
	"math/rand"

	"github.com/PavelDonchenko/sensor-go/internal/domain"
)

func GenerateCoordinates(id int) domain.Coordinates {
	minOffset := float64(id) * -0.1 // Minimum offset value
	maxOffset := float64(id) * 0.1  // Maximum offset value

	// Generate random offset values for X, Y, Z coordinates within the specified range
	offsetX := minOffset + rand.Float64()*(maxOffset-minOffset)
	offsetY := minOffset + rand.Float64()*(maxOffset-minOffset)
	offsetZ := minOffset + rand.Float64()*(maxOffset-minOffset)

	// Generate a random base coordinate within a specific range
	baseX := -10.0 + rand.Float64()*(10.0-(-10.0))
	baseY := -10.0 + rand.Float64()*(10.0-(-10.0))
	baseZ := -10.0 + rand.Float64()*(0.0-(-10.0))

	// Calculate the final coordinates by applying the offsets to the base coordinate
	x := baseX + offsetX
	y := baseY + offsetY
	z := baseZ + offsetZ

	return domain.Coordinates{
		X: x,
		Y: y,
		Z: z,
	}
}
