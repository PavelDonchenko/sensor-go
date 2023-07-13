package generations

import "math/rand"

func GenerateTransparency(previousTransparency int) int {
	minTransparency := previousTransparency - 10
	if minTransparency < 0 {
		minTransparency = 0
	}
	maxTransparency := previousTransparency + 10
	if maxTransparency > 100 {
		maxTransparency = 100
	}

	return rand.Intn(maxTransparency-minTransparency+1) + minTransparency
}
