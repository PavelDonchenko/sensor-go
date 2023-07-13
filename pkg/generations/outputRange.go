package generations

import (
	"math/rand"
)

func GenerateRandomInt() int {
	values := []int{5, 10, 15, 20, 25}

	randomIndex := rand.Intn(len(values))

	randomInt := values[randomIndex]

	return randomInt
}
