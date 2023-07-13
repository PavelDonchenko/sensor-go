package generations

import (
	"math"
	"math/rand"
)

func GenerateTemperature(depth float64) float64 {
	randomFloat := rand.Float64() + 1

	scaledFloat := randomFloat * 2

	result := scaledFloat + math.Abs(depth)*3

	roundedResult := math.Round(result*1000) / 1000

	return roundedResult
}
